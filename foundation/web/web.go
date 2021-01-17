package web

import (
    "context"
    "syscall"
    "time"
    "os"
    "net/http"
    "github.com/dimfeld/httptreemux/v5"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
)

type ctxKey int

const KeyValues ctxKey = 1

type Values struct {
    TraceID string
    Now time.Time
    StatusCode int
}


type App struct {
    mux *httptreemux.ContextMux
    otmux http.Handler
    shutdown chan os.Signal
    mw []Middleware
}

type Handler  func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

func NewApp(shutdown chan os.Signal, mw ...Middleware) *App {
    mux := httptreemux.NewContextMux()
    return &App{
        mux:mux,
    // https://w3c.github.io/trace-context/
        otmux: otelhttp.NewHandler(mux, "request"),
        shutdown: shutdown,
        mw: mw,
    }

}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    a.otmux.ServeHTTP(w, r)
}

func (a *App) Handle(method string, path string, handler Handler, mw ...Middleware) {

    handler = wrapMiddleware(mw, handler)

    handler = wrapMiddleware(a.mw, handler)

    h := func(w http.ResponseWriter, r *http.Request) {
        //BOILERPLAE

        ctx := r.Context()
		ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, r.URL.Path)
        defer span.End()

        v := Values{
            TraceID: span.SpanContext().TraceID.String(),
            Now: time.Now(),
        }

        ctx = context.WithValue(ctx, KeyValues, &v)


        if err := handler(ctx, w, r); err != nil {
            a.SignalShutdown()
            return
        }

        //BOILERPLAE
    }
    a.mux.Handle(method, path, h)
}


//SignalShutdown is used to gracefully shutdown the app when an intergrity
// issue is identified

func (a *App)SignalShutdown() {
    a.shutdown <- syscall.SIGTERM
}
