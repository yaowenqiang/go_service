package handlers
import (
    "log"
    "net/http"
    "encoding/json"
)


type check struct {
    Log *log.Logger
}


func (c check) readiness(w http.ResponseWriter, r *http.Request)  {
        status := struct {
            Status string
        }{
            Status: "OK",
        }
        json.NewEncoder(w).Encode(status)
        log.Println(r, status)
}

