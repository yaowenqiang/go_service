package main
import (
    "crypto/rsa"
    "log"
    "os"
    "time"
    "os/signal"
    "context"
    "syscall"
    "fmt"
    "io/ioutil"
    "expvar"
    "net/http"
    _ "net/http/pprof"

    "github.com/pkg/errors"
    "github.com/dimfeld/httptreemux/v5"
	"github.com/ardanlabs/conf"
	"github.com/yaowenqiang/service/app/sales-api/handlers"
	"github.com/yaowenqiang/service/business/auth"
	"github.com/yaowenqiang/service/foundation/database"
    "github.com/dgrijalva/jwt-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/trace/zipkin"
	"go.opentelemetry.io/otel/sdk/trace"
)

/*
Need to figure out timeouts for http service.
You might want to reset your DB_HOST env var during test tear down.
Service should start even without a DB running yet.
symbols in profiles: https://github.com/golang/go/issues/23376 / https://github.com/google/pprof/pull/366
*/


// build is the git version of this program. It is set using build flags in the makefile.
var build = "develop"


func main() {
    log := log.New(os.Stdout, "SALES: ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
    if err := run(log); err != nil {
        log.Println("main error: ", err)
        os.Exit(1)

        //log.Fatalln()
    }
}


func run(log *log.Logger) error{
	var cfg struct {
		conf.Version
		Web struct {
			APIHost         string        `conf:"default:0.0.0.0:3000"`
			DebugHost       string        `conf:"default:0.0.0.0:4000"`
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:5s,noprint"`
			ShutdownTimeout time.Duration `conf:"default:5s"`
		}
		DB struct {
			User         	string        `conf:"default:postgres"`
			Password        string        `conf:"default:postgres, noprint"`
			Host         	string        `conf:"default:localhost"`
			Name         	string        `conf:"default:postgres"`
			DisableTLS      bool          `conf:"default:true"`
		}
		Zipkin struct {
		ReportURI string `conf:"default:http://0.0.0.0:9411/api/v2/spans"`
		ServiceName string `conf:"default:sales-api"`
		Probability float64 `conf:"default:0.05"`
		}
		Auth struct {
			KeyID         		string        `conf:"default:asdlfjldasjfdsjfldasjfl jlsjflweqjio;ewjejf"`
			PrivateKeyFile      string        `conf:"default:private.pem"`
			Algorithm     		string 		  `conf:"default:RS256"`
		}
	}


	cfg.Version.SVN = build
	cfg.Version.Desc = "copyright information here"

	if err := conf.Parse(os.Args[1:], "SALES", &cfg); err != nil {
		if err := conf.Parse(os.Args[1:], "SALES", &cfg); err != nil {
			switch err {
			case conf.ErrHelpWanted:
				usage, err := conf.Usage("SALES", &cfg)
				if err != nil {
					return errors.Wrap(err, "generating config usage")
				}
				fmt.Println(usage)
				return nil
			case conf.ErrVersionWanted:
				version, err := conf.VersionString("SALES", &cfg)
				if err != nil {
					return errors.Wrap(err, "generating config version")
				}
				fmt.Println(version)
				return nil
			}
			return errors.Wrap(err, "parsing config")
		}

	}

	// =========================================================================
	// App Starting

	// Print the build version for our logs. Also expose it under /debug/vars.
	expvar.NewString("build").Set(build)
	log.Printf("main : Started : Application initializing : version %q", build)
	defer log.Println("main: Completed")

	out, err := conf.String(&cfg)
	if err != nil {
		return errors.Wrap(err, "generating config for output")
	}
	log.Printf("main: Config :\n%v\n", out)



	// =========================================================================
	// Initializing authentication support

	log.Println("main: Started : Intializing authentication support")

	privatePEM, err := ioutil.ReadFile(cfg.Auth.PrivateKeyFile)

	if err != nil {
		return errors.Wrap(err, "reading auth private key")
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privatePEM)


	if err != nil {
		return errors.Wrap(err, "parsing auth private key")
	}

	lookup := func(kid string) (*rsa.PublicKey, error) {
		switch kid {
		case cfg.Auth.KeyID:
			return &privateKey.PublicKey, nil
		}

		return nil, fmt.Errorf("no public key found for the specified kid: %s", kid)
	}

	auth, err  := auth.New(cfg.Auth.Algorithm, lookup, auth.Keys{cfg.Auth.KeyID: privateKey})

	if err != nil {
		return errors.Wrap(err, "constructing auth")
	}
	// =========================================================================
	// Start Database
	log.Println("main: Intializing database support")
	db, err := database.Open(database.Config{
		User: cfg.DB.User,
		Password: cfg.DB.Password,
		Host: cfg.DB.Host,
		Name: cfg.DB.Name,
		DisableTLS: cfg.DB.DisableTLS,
	})

	// Start Tracing Support

	log.Println("main: Intialization OT/Zipkin tracing support")

	exporter, err := zipkin.NewRawExporter(
		cfg.Zipkin.ReportURI,
		cfg.Zipkin.ServiceName,
		zipkin.WithLogger(log),
	)

	if err != nil {
		return errors.Wrap(err, "creating new exporter")
	}
	tp := trace.NewTracerProvider(
		trace.WithConfig(trace.Config{DefaultSampler: trace.TraceIDRatioBased(cfg.Zipkin.Probability)}),
		trace.WithBatcher(exporter,
			trace.WithMaxExportBatchSize(trace.DefaultMaxExportBatchSize),
			trace.WithBatchTimeout(trace.DefaultBatchTimeout),
			trace.WithMaxExportBatchSize(trace.DefaultMaxExportBatchSize),
		),
	)

	otel.SetTracerProvider(tp)

	if err != nil {
		return errors.Wrap(err, "connecting to db")
	}

	defer func() {
		log.Printf("main: database Stopping: %s", cfg.DB.Host)
		db.Close()
	}()



	// =========================================================================
	// start debug service
	//
	// /debug/pprof - added to the default mux by importing the net/http/pprof package.
	// /debug/vars - added to the default mux by importing the expvar package.
	//
	// not concerned with shutting this down when the application is shutdown.

	log.Println("main: Initializing debugging support")

	go func() {
		log.Printf("main: Debug Listening %s", cfg.Web.DebugHost)
		if err := http.ListenAndServe(cfg.Web.DebugHost, http.DefaultServeMux); err != nil {
			log.Printf("main: Debug Listener closed : %v", err)
		}
	}()


		m := httptreemux.NewContextMux()
		m.Handle(http.MethodGet, "/test", nil)

// =========================================================================
	// Start API Service

	log.Println("main: Initializing API support")

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	api := http.Server{
		Addr:         cfg.Web.APIHost,
		Handler:      handlers.API(build, shutdown, log, auth, db),
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
	}


	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	// Start the service listening for requests.
	go func() {
		log.Printf("main: API listening on %s", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

// =========================================================================
	// Shutdown

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		return errors.Wrap(err, "server error")

	case sig := <-shutdown:
		log.Printf("main: %v : Start shutdown", sig)

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		// Asking listener to shutdown and shed load.
		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return errors.Wrap(err, "could not stop server gracefully")
		}
	}

	return nil
}
