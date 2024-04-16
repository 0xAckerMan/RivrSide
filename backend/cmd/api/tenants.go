package main

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/0xAckerMan/internal/data"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (app *Application) HandleUpdateTenantInfo(w http.ResponseWriter, r *http.Request) {
	currentTenant := r.Context().Value("user").(data.User)

	var input struct {
		First_name   *string `json:"first_name"`
		Last_name    *string `json:"last_name"`
		Email        *string `json:"email"`
		PhoneNumber  *string `json:"phone_number"`
		Password     *string `json:"password"`
		Organisation *string `json:"organisation"`
		Position     *string `json:"position"`
		Subcription  *int    `json:"subscription"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	res := app.DB.First(&currentTenant, currentTenant.ID)
	if res.RowsAffected == 0 {
		app.noRecordFoundResponse(w, r)
		return
	}

	if res.Error != nil {
		app.serverErrorResponse(w, r, res.Error)
		return
	}

	if input.First_name != nil {
		currentTenant.First_name = *input.First_name
	}

	if input.Last_name != nil {
		currentTenant.Last_name = *input.Last_name
	}

	if input.Email != nil {
		currentTenant.Email = *input.Email
	}

	if input.PhoneNumber != nil {
		currentTenant.PhoneNumber = *input.PhoneNumber
	}

	if input.Organisation != nil {
		currentTenant.Organisation = *input.Organisation
	}

	if input.Position != nil {
		currentTenant.Position = *input.Position
	}

	if input.Password != nil {
		hash, err := bcrypt.GenerateFromPassword([]byte(*input.Password), 12)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		currentTenant.Password = string(hash)
	}
	var subscription data.Subscription
	if input.Subcription != nil {
		if err := app.DB.Unscoped().First(&subscription, "user_id = ?", currentTenant.ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				subscription = data.Subscription{
					UserID:    int(currentTenant.ID),
					PackageID: *input.Subcription,
				}

				if err = app.DB.Create(&subscription).Error; err != nil {
					app.serverErrorResponse(w, r, err)
					return
				}
			} else {
				app.serverErrorResponse(w, r, err)
				return
			}
		} else {

			subscription.PackageID = *input.Subcription
			if err = app.DB.Save(&subscription).Error; err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}

			if subscription.DeletedAt.Valid {
				if err = app.DB.Unscoped().Model(&subscription).Update("deleted_at", nil).Error; err != nil {
					app.serverErrorResponse(w, r, err)
					return
				}
			}
		}
	}

	res = app.DB.Save(&currentTenant)
	if res.RowsAffected == 0 {
		app.noRecordFoundResponse(w, r)
		return
	}

	if res.Error != nil {
		app.serverErrorResponse(w, r, res.Error)
		return
	}

	// recheck

	if err = app.DB.Model(&currentTenant).Preload("Room").Preload("Role").First(&currentTenant).Error; err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err = app.DB.Model(&subscription).Preload("PackagePlan").First(&subscription).Error; err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	updatedUser := data.TenantInfo{
		User:      currentTenant,
		PackageID: int64(subscription.PackageID),
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"tenant": updatedUser}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *Application) HandleGetRoomInfo(w http.ResponseWriter, r *http.Request) {
	return
}

func (app *Application) HandleGetProfileInfo(w http.ResponseWriter, r *http.Request) {
	currentTenant := r.Context().Value("user").(data.User)
	res := app.DB.Preload("Role").Preload("Room").First(&currentTenant, currentTenant.ID)
	if res.RowsAffected == 0 {
		app.noRecordFoundResponse(w, r)
		return
	}

	if res.Error != nil {
		app.serverErrorResponse(w, r, res.Error)
		return
	}

	err := app.writeJSON(w, http.StatusOK, envelope{"my_profile": currentTenant}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *Application) HandleGetMySubscription(w http.ResponseWriter, r *http.Request) {
	return
}

func (app *Application) HandleRequestTransfer(w http.ResponseWriter, r http.Request) {
	return
}

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
