package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var Version = "1.0.0"

type Config struct {
	env  string
	port int
}

type Application struct {
	Config
	logger *log.Logger
}

func main() {
	var cfg Config

	flag.StringVar(&cfg.env, "env", "env", "The environment of the api")
	flag.IntVar(&cfg.port, "port", 3000, "The running port")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	app := &Application{
		Config: cfg,
		logger: logger,
	}

    addr := fmt.Sprintf(":%d", cfg.port)

	logger.Printf("%s server running on port %d", app.env, app.port)

    srv := http.Server{
        Addr: addr,
        Handler: app.routes(),
        IdleTimeout: time.Minute,
        ReadTimeout: 30 * time.Second,
        WriteTimeout: 30 * time.Second,
    }
    srv.ListenAndServe()
}
