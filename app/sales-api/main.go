package main
import (
    "log"
    "os"
    "time"
    "fmt"
    "expvar"
    "net/http"
    _ "net/http/pprof"

    "github.com/pkg/errors"
    "github.com/dimfeld/httptreemux/v5"
	"github.com/ardanlabs/conf"
)

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
	// Start Debug Service
	//
	// /debug/pprof - Added to the default mux by importing the net/http/pprof package.
	// /debug/vars - Added to the default mux by importing the expvar package.
	//
	// Not concerned with shutting this down when the application is shutdown.

	log.Println("main: Initializing debugging support")

	go func() {
		log.Printf("main: Debug Listening %s", cfg.Web.DebugHost)
		if err := http.ListenAndServe(cfg.Web.DebugHost, http.DefaultServeMux); err != nil {
			log.Printf("main: Debug Listener closed : %v", err)
		}
	}()

	select {}

		m := httptreemux.NewContextMux()
		m.Handle(http.MethodGet, "/test", nil)
		return nil
}
