package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Version = "1.0.0"

type Config struct {
	env  string
	port int
    db struct{
        dsn string
    }
}

type Application struct {
	Config Config
	logger *log.Logger
    DB *gorm.DB
}

func init(){
    LoadEnv()
}

func main() {
	var cfg Config

	flag.StringVar(&cfg.env, "env", "dev", "The environment of the api")
	flag.IntVar(&cfg.port, "port", 3000, "The running port")
    flag.StringVar(&cfg.db.dsn, "dsn", os.Getenv("DATABASE_DSN"), "The db connection dsn")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := gorm.Open(postgres.Open(cfg.db.dsn), &gorm.Config{})
	if err != nil {
		logger.Fatal(err)
	}

	DB, err := db.DB()
	if err != nil {
		logger.Fatal(err)
	}
	defer DB.Close()
    logger.Printf("database connection pool established")

	app := &Application{
        DB: db,
		Config: cfg,
		logger: logger,
	}
    app.migrations()

    addr := fmt.Sprintf(":%d", cfg.port)

	logger.Printf("%s server running on port %d", app.Config.env, app.Config.port)

    srv := http.Server{
        Addr: addr,
        Handler: app.routes(),
        IdleTimeout: time.Minute,
        ReadTimeout: 30 * time.Second,
        WriteTimeout: 30 * time.Second,
    }
    srv.ListenAndServe()
}
