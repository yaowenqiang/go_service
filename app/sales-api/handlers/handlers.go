// Package handlers contains the full set of handler functions and routes
// supported by the web api.
package handlers

import (
    "os"
    "log"
    "net/http"
    "github.com/dimfeld/httptreemux"
)


// API constructs an http.Handler with all application routes defined.
func API(build string, shutdown chan os.Signal, log *log.Logger)  *httptreemux.ContextMux {

    tm := httptreemux.NewContextMux()
    c := check{
        Log: log,
    }

    tm.Handle(http.MethodGet, "/test",c.readiness)
    return tm
}

