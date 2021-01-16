// Package handlers contains the full set of handler functions and routes
// supported by the web api.
package handlers

import (
    "os"
    "log"
    "net/http"
    "github.com/yaowenqiang/service/business/mid"
    "github.com/yaowenqiang/service/foundation/web"
    "github.com/yaowenqiang/service/business/data/user"
    "github.com/yaowenqiang/service/business/auth"
    "github.com/jmoiron/sqlx"
)


// API constructs an http.Handler with all application routes defined.
func API(build string, shutdown chan os.Signal, log *log.Logger, a *auth.Auth, db *sqlx.DB) *web.App {

    app := web.NewApp(shutdown, mid.Logger(log), mid.Errors(log), mid.Metrics(), mid.Panic(log))
    cg := checkGroup{
        build: build,
        Log: log,
		db: db,
    }

    //app.Handle(http.MethodGet, "/readiness",c.readiness, mid.Authenticate(a), mid.Authorize(log, auth.RoleAdmin))
    //app.Handle(http.MethodGet, "/liveness",c.readiness, mid.Authenticate(a), mid.Authorize(log, auth.RoleAdmin))
    app.Handle(http.MethodGet, "/readiness",cg.readiness)
    app.Handle(http.MethodGet, "/liveness",cg.liveness)


	ug := userGroup{
		user: user.New(log, db),
		auth: a,
	}
	app.Handle(http.MethodGet, "/users/:page/:rows",ug.query, mid.Authenticate(a), mid.Authorize(auth.RoleAdmin))
	app.Handle(http.MethodGet, "/users/:id",ug.queryByID, mid.Authenticate(a))
	app.Handle(http.MethodGet, "/users/token/:kid",ug.token)
	app.Handle(http.MethodPost, "/users/",ug.create, mid.Authenticate(a), mid.Authorize(auth.RoleAdmin))
	app.Handle(http.MethodPut, "/users/:id",ug.update, mid.Authenticate(a), mid.Authorize(auth.RoleAdmin))
	app.Handle(http.MethodDelete, "/users/:id",ug.delete, mid.Authenticate(a), mid.Authorize(auth.RoleAdmin))

    return app
}

