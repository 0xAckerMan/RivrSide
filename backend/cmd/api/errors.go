package main

import (
	"fmt"
	"net/http"
)

func (app *Application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}){
    app.writeJSON(w,status,envelope{"response": message},nil)
}

func (app *Application) logError(err error){
    app.logger.Println(err)
}

func (app *Application) notFoundResponse(w http.ResponseWriter, r *http.Request){
    message := "The requested resource couldnt be found at the moment"
    app.errorResponse(w,r,http.StatusNotFound,message)
}

func (app *Application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request){
    message := fmt.Sprintf("%s is not allowed for this operation", r.Method)
    app.errorResponse(w,r,http.StatusMethodNotAllowed,message)
}

func (app *Application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error){
    app.logError(err)
    message := "Sorry, the server is experiencing an error at the moment"
    app.errorResponse(w,r,http.StatusInternalServerError,message)
}
