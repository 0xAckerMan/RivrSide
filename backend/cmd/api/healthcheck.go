package main

import (
	"net/http"
)

func (app *Application) healthcheck (w http.ResponseWriter, r *http.Request) {
    status := map[string]string{
        "status": "active",
        "environment": app.Config.env,
        "version": Version,
    }
    err := app.writeJSON(w,http.StatusOK, envelope{"health":status},nil)
    if err != nil{
        app.logger.Println(err)
        return
    }
}

func (app *Application) ApiHome (w http.ResponseWriter, r *http.Request){
    message := "Welcome!!! You made it. Your journey begins here"
    err := app.writeJSON(w,http.StatusOK,envelope{"success": message},nil)
    if err != nil{
        app.serverErrorResponse(w,r,err)
        return
    }
}
