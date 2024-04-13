package main

import (
	"net/http"

	"github.com/0xAckerMan/internal/data"
	"golang.org/x/crypto/bcrypt"
)

func (app *Application) HandleUpdateTenantInfo(w http.ResponseWriter, r *http.Request) {
	currentTenant := r.Context().Value("tenant").(data.User)

	var input struct {
		First_name   *string `json:"first_name"`
		Last_name    *string `json:"last_name"`
		Email        *string `json:"email"`
		PhoneNumber  *string `json:"phone_number"`
		Password     *string `json:"password"`
		Organisation *string `json:"organisation"`
		Position     *string `json:"position"`
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

	err = app.writeJSON(w, http.StatusOK, envelope{"tenant": currentTenant}, nil)
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
