package handlers
import (
    "log"
    "math/rand"
    "context"
    "net/http"
    "github.com/pkg/errors"
    "github.com/yaowenqiang/service/foundation/web"
)


type check struct {
    Log *log.Logger
}


//func (c check) readiness(ctx context.Context, w http.ResponseWriter, r *http.Request) error  {
func (c check) readiness(ctx context.Context,  w http.ResponseWriter, r *http.Request)  error {
    if m := rand.Intn(100); m % 2 == 0 {
        return errors.New("untrusted error")
    }
        status := struct {
            Status string
        }{
            Status: "OK",
        }
        return web.Respond(ctx, w, status, http.StatusOK)
}

