package main

import (
	"errors"
	"fmt"
	"net/http"
)

var ErrRecordNotFound = errors.New("record searched could not be found")

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

func (app *Application) noRecordFoundResponse(w http.ResponseWriter, r *http.Request){
    message := "no record found"
    app.errorResponse(w,r,http.StatusNotFound,message)
}

func (app *Application) duplicateRecordResponse(w http.ResponseWriter, r *http.Request){
    message:= "duplicate input data, record already exists"
    app.errorResponse(w,r,http.StatusOK,message)
}

