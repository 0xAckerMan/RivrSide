package main

import (
	"net/http"
)

func (app *Application) healthcheck (w http.ResponseWriter, r *http.Request) {
    status := map[string]string{
        "status": "active",
        "environment": app.env,
        "version": Version,
    }
    err := app.writeJSON(w,http.StatusOK, envelope{"health":status},nil)
    if err != nil{
        app.logger.Println(err)
        return
    }
}
