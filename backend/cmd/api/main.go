package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
)

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

  logger.Printf("%s server running on port %d", app.env,app.port)

  http.HandleFunc("/healthcheck",app.healthcheck)
	http.ListenAndServe(":3000", nil)
}

func (app *Application) healthcheck(w http.ResponseWriter, r *http.Request){
  status := map[string]interface{}{
    "status": "active",
    "env": app.env,
    "version": "1.0.0",
  }

  js, err := json.Marshal(status)
  if err != nil{
    app.logger.Println("cant unmarshal the data")
    http.Error(w,"Internal server error", http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.Write(js)
}
