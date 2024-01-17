package main

import (
	"encoding/json"
	"net/http"
)

func (app *Application) healthcheck (w http.ResponseWriter, r *http.Request){
    status := map[string]string{
        "status": "active",
        "environment": app.env,
        "version": Version,
    }

    js, err := json.Marshal(status)
    if err != nil{
        http.Error(w, "internal server error", http.StatusInternalServerError)
        return
    }

    js = append(js, '\n')

    w.Header().Set("Content-Type", "application/json")

    w.Write(js)
}
