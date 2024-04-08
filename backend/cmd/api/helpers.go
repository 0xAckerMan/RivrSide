package main

import (
	"encoding/json"
	"net/http"
)

type envelope map[string]interface{}

func (app *Application) writeJSON(w http.ResponseWriter, status int, message envelope, header http.Header) error {
    js, err := json.Marshal(message)
    if err != nil{
        return err
    }

    js =append(js, '\n')

    for key, value := range header{
        w.Header()[key] = value
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    w.Write(js)

    return err
}
