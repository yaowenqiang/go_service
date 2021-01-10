// Package handlers contains the full set of handler functions and routes
// supported by the web api.
package handlers

import (
    "os"
    "log"
    "net/http"
    "github.com/yaowenqiang/service/business/mid"
    "github.com/yaowenqiang/service/foundation/web"
    "github.com/yaowenqiang/service/business/auth"
)


// API constructs an http.Handler with all application routes defined.
func API(build string, shutdown chan os.Signal, log *log.Logger, a *auth.Auth) *web.App {

    app := web.NewApp(shutdown, mid.Logger(log), mid.Errors(log), mid.Metrics(), mid.Panic(log))
    c := check{
        build: build,
        Log: log,
    }

    //app.Handle(http.MethodGet, "/readiness",c.readiness, mid.Authenticate(a), mid.Authorize(log, auth.RoleAdmin))
    //app.Handle(http.MethodGet, "/liveness",c.readiness, mid.Authenticate(a), mid.Authorize(log, auth.RoleAdmin))
    app.Handle(http.MethodGet, "/readiness",c.readiness)
    app.Handle(http.MethodGet, "/liveness",c.liveness)
    return app
}

