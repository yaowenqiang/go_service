package main

import (
    "log"
    "os"

    "github.com/pkg/errors"
)

func main() {
    if err := run(); err != nil {
        log.Println(err)
        os.Exit(1)
    }
}


func run() error{
    return errors.New("random error")
}
