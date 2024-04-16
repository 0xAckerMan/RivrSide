package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/0xAckerMan/internal/data"
	"gorm.io/gorm"
)

func (app *Application) HandleMakePayment(w http.ResponseWriter, r *http.Request) {
	currentTenant := r.Context().Value("user").(data.User)

	var input data.MakePayment
	var subscription data.Subscription

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	res := app.DB.Preload("User").Preload("PackagePlan").Where("user_id = ? AND deleted_at IS NULL", currentTenant.ID).
		Order("id").Limit(1).First(&subscription)
	if res.Error == gorm.ErrRecordNotFound {
		message := "You do not have a package, please select before making payments"
		app.errorResponse(w, r, http.StatusOK, envelope{"response": message})
		return
	}

	if res.Error != nil {
		app.serverErrorResponse(w, r, res.Error)
		return
	}

	var payment data.Payment

	payment = data.Payment{
		UserID:         int64(currentTenant.ID),
		Amount:         input.Amount,
		SubscriptionID: int(subscription.ID),
		Month:          time.Now().Month(),
		Year:           time.Now().Year(),
		Subscription:   &subscription,
	}

	var totalPaid sql.NullFloat64

    result := app.DB.Model(&data.Payment{}).
        Where("user_id = ? AND ((month = ? AND year = ?) OR (year = ? AND month < ?))",
            currentTenant.ID, time.Now().Month(), time.Now().Year(), time.Now().Year(), time.Now().Month()).
        Pluck("SUM(amount)", &totalPaid)

	if result.Error != nil {
		app.serverErrorResponse(w, r, result.Error)
		return
	}
	var totalPaidAmount float64
	if totalPaid.Valid {
		totalPaidAmount = totalPaid.Float64
	}
	totalPaidAmount += float64(input.Amount)

	payment.Balance = subscription.PackagePlan.Price - int(totalPaidAmount)
	if int(totalPaidAmount) > payment.Subscription.PackagePlan.Price || payment.Balance < 0 {
		nextMonth := data.Payment{
			UserID:         int64(currentTenant.ID),
			Amount:         -payment.Amount,
			SubscriptionID: payment.SubscriptionID,
			Month:          time.Now().AddDate(0, 1, 0).Month(),
		}
		if payment.Year == int(time.December) {
			nextMonth.Year = time.Now().AddDate(1, 0, 0).Year()
		} else {
			nextMonth.Year = time.Now().Year()
		}

		if err = app.DB.Create(&nextMonth).Error; err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	if err = app.DB.Create(&payment).Error; err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	res = app.DB.Model(&payment).Preload("User").Preload("Subscription").First(&payment)
	if res.Error != nil {
		app.serverErrorResponse(w, r, res.Error)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"payment": payment}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

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
