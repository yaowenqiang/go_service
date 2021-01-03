// Package handlers contains the full set of handler functions and routes
// supported by the web api.
package handlers

import (
    "os"
    "log"
    "net/http"
    "github.com/yaowenqiang/service/business/mid"
    "github.com/yaowenqiang/service/foundation/web"
)


// API constructs an http.Handler with all application routes defined.
func API(build string, shutdown chan os.Signal, log *log.Logger)  *web.App {

    app := web.NewApp(shutdown, mid.Logger(log), mid.Errors(log))
    c := check{
        Log: log,
    }

    app.Handle(http.MethodGet, "/readiness",c.readiness)
    return app
}

