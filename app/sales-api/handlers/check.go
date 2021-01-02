package handlers
import (
    "log"
    "context"
    "net/http"
    "encoding/json"
)


type check struct {
    Log *log.Logger
}


//func (c check) readiness(ctx context.Context, w http.ResponseWriter, r *http.Request) error  {
func (c check) readiness(ctx context.Context,  w http.ResponseWriter, r *http.Request)  error {
        status := struct {
            Status string
        }{
            Status: "OK",
        }
        log.Println(r, status)
        return json.NewEncoder(w).Encode(status)
}

