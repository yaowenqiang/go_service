package handlers
import (
    "log"
    _ "math/rand"
    "context"
    "net/http"
    _ "github.com/pkg/errors"
    "github.com/yaowenqiang/service/foundation/web"
)


type check struct {
    Log *log.Logger
}


//func (c check) readiness(ctx context.Context, w http.ResponseWriter, r *http.Request) error  {
func (c check) readiness(ctx context.Context,  w http.ResponseWriter, r *http.Request)  error {
    //if m := rand.Intn(100); m % 2 == 0 {
        //return errors.New("untrusted error")
    //} else {
        //return web.NewRequestError(errors.New("tusted error"), http.StatusNotFound)
        //return web.NewRequestError(errors.New("tusted error"), http.StatusNotFound)
        //panic("forcing panic")
        //return web.NewShutdownError("forcing shutdown")

    //}
        status := struct {
            Status string
        }{
            Status: "OK",
        }
        return web.Respond(ctx, w, status, http.StatusOK)
}

