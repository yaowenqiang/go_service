package web

import (
    "context"
    "syscall"
    "time"
    "os"
    "net/http"
    "github.com/dimfeld/httptreemux/v5"
    "github.com/google/uuid"
)

type ctxKey int

const KeyValues ctxKey = 1

type Values struct {
    TraceID string
    Now time.Time
    StatusCode int
}


type App struct {
    *httptreemux.ContextMux
    shutdown chan os.Signal
    mw []Middleware
}

type Handler  func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

func NewApp(shutdown chan os.Signal, mw ...Middleware) *App {
    app := App{
        ContextMux: httptreemux.NewContextMux(),
        shutdown: shutdown,
        mw: mw,
    }
    return &app
}

func (a *App) Handle(method string, path string, handler Handler, mw ...Middleware) {

    handler = wrapMiddleware(mw, handler)

    handler = wrapMiddleware(a.mw, handler)

    h := func(w http.ResponseWriter, r *http.Request) {
        //BOILERPLAE

        v := Values{
            TraceID: uuid.New().String(),
            Now: time.Now(),
        }

        ctx := context.WithValue(r.Context(), KeyValues, &v)


        if err := handler(ctx, w, r); err != nil {
            a.SignalShutdown()
            return
        }

        //BOILERPLAE
    }
    a.ContextMux.Handle(method, path, h)
}


//SignalShutdown is used to gracefully shutdown the app when an intergrity
// issue is identified

func (a *App)SignalShutdown() {
    a.shutdown <- syscall.SIGTERM
}
