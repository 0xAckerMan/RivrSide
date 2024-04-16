package main

import (
	"net/http"

	"github.com/0xAckerMan/internal/data"
)

func (app *Application) HandleGetAllPayment  (w http.ResponseWriter, r *http.Request){
    var payment []data.Payment

    result := app.DB.Preload("User").Preload("Subscription").Find(&payment)
    if result.RowsAffected ==0{
        app.noRecordFoundResponse(w,r)
        return
    }

    if result.Error != nil{
        app.serverErrorResponse(w,r,result.Error)
        return
    }

    err := app.writeJSON(w,http.StatusOK,envelope{"all_payments": payment}, nil)
    if err != nil{
        app.serverErrorResponse(w,r,err)
        return
    }
}

func (app *Application) HandleGetAllPaymentPerMonth (w http.ResponseWriter, r *http.Request){
    return
}

func (app *Application) HandleGetCurrentMonthPayment (w http.ResponseWriter, r *http.Request){
    return
}

func (app *Application) HandleAllCurrentMonthMealandRoom (w http.ResponseWriter, r *http.Request){
    return
}

func (app *Application) HandleAllCurrentMonthRoomOnly (w http.ResponseWriter, r *http.Request){
    return
}

func (app *Application) HandleGetSingleUserPaymentHistory (w http.ResponseWriter, r *http.Request){
    return
}

func (app *Application) HandleUserBalance (w http.ResponseWriter, r *http.Request){
    return
}
