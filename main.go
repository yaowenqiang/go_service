package main

import (
    "log"
    "os"
    "net/http"

    "github.com/pkg/errors"
    "github.com/dimfeld/httptreemux/v5"
)

func main() {
    if err := run(); err != nil {
        log.Println(err)
        os.Exit(1)
    }
}


func run() error{
    m := httptreemux.NewContextMux()
    m.Handle(http.MethodGet, "/test", nil)
    return errors.New("random error")
}
