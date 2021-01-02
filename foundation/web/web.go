package web

import (
    "context"
    "syscall"
    "os"
    "net/http"
    "github.com/dimfeld/httptreemux/v5"
)

type App struct {
    *httptreemux.ContextMux
    shutdown chan os.Signal
}

type Handler  func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

func NewApp(shutdown chan os.Signal) *App {
    app := App{
        ContextMux: httptreemux.NewContextMux(),
        shutdown: shutdown,
    }
    return &app
}

func (a *App) Handle(method string, path string, handler Handler) {
    h := func(w http.ResponseWriter, r *http.Request) {
        //BOILERPLAE

        if err := handler(r.Context(), w, r); err != nil {
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
