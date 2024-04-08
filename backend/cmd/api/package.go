package main

import (
	"net/http"

	"github.com/0xAckerMan/internal/data"
)

func (app *Application) handle_getAllPackages(w http.ResponseWriter, r *http.Request){
    err := app.writeJSON(w,http.StatusOK, envelope{"packages": data.NewPackage()},nil)
    if err != nil{
        app.logger.Println(err)
        return
    }
}

func (app *Application) handle_getSinglePackage(w http.ResponseWriter, r *http.Request){
    err := app.writeJSON(w,200,envelope{"message":"Single package"},nil)
    if err != nil{
        app.logger.Fatal(err)
        return
    }
}
