package main
import (
    "log"
    "os"
)


func main() {
    log := log.New(os.Stdout, "SALES: ", log.LstdFlags|logs.Lmicroseconds|log.Lshortfile)
    if err := run(log); err != nil {
        log.Println("main error: ", err)
        os.Exit(1)

        //log.Fatalln()
    }
}


func run(log *log.Logger) error {
    return nil
}
