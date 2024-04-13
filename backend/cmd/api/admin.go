package main

import (
	"net/http"

	"github.com/0xAckerMan/internal/data"
	"golang.org/x/crypto/bcrypt"
)

func (app *Application) HandleCreateManager(w http.ResponseWriter, r *http.Request) {
	var input data.CreateManager

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	manager := data.User{
		First_name:   input.First_name,
		Last_name:    input.Last_name,
		Email:        input.Email,
		PhoneNumber:  input.PhoneNumber,
		Password:     input.Password,
		Gender:       input.Gender,
		RoomID:       0,
		Position:     "manager",
		Organisation: "riverside",
		Ismanager:    true,
		IsActive:     true,
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), 12)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	manager.Password = string(hash)

	var role data.Role
	roleResult := app.DB.Where(&data.Role{Role: "manager"}).First(&role)
	if roleResult.RowsAffected == 0 {
		message := "Create a role named manager first please"
		app.errorResponse(w, r, http.StatusInternalServerError, envelope{"error": message})
		return
	}

	if roleResult.Error != nil {
		app.serverErrorResponse(w, r, roleResult.Error)
		return
	}

	manager.RoleID = int64(role.ID)

	res := app.DB.Create(&manager)
	if res.Error != nil {
		app.serverErrorResponse(w, r, res.Error)
		return
	}

	if err := app.DB.Model(&manager).Preload("Role").First(&manager).Error; err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err = app.writeJSON(w, http.StatusCreated, envelope{"manager": manager}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *Application) HandleGetAllManagers(w http.ResponseWriter, r *http.Request) {
	var role data.Role

	res := app.DB.Where(data.Role{Role: "manager"}).First(&role)
	if res.RowsAffected == 0 {
		message := "Oops the role manager does not exist, please create it first"
		app.errorResponse(w, r, http.StatusNotFound, envelope{"response": message})
		return
	}

	if res.Error != nil {
		app.serverErrorResponse(w, r, res.Error)
		return
	}

	var managerusers []data.User
	res = app.DB.Preload("Role").Where(data.User{RoleID: int64(role.ID)}).Find(&managerusers)
	if res.RowsAffected == 0 {
		app.notFoundResponse(w, r)
		return
	}

	if res.Error != nil {
		app.serverErrorResponse(w, r, res.Error)
		return
	}

	err := app.writeJSON(w, http.StatusOK, envelope{"manager": managerusers}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *Application) HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDparam(w, r)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	var user data.User

	res := app.DB.Delete(&user, id)
	if res.RowsAffected == 0 {
		app.noRecordFoundResponse(w, r)
		return
	}

	if res.Error != nil {
		app.serverErrorResponse(w, r, res.Error)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"success": "record deleted successfully"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *Application) HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDparam(w, r)
	if err != nil {
		app.serverErrorResponse(w, r,err)
		return
	}

	var input struct {
		First_name   *string `json:"first_name"`
		Last_name    *string `json:"last_name"`
		Email        *string `json:"email"`
		PhoneNumber  *string `json:"phone_number"`
		Password     *string `json:"password"`
		Gender       *string `json:"gender"`
		RoomID       *int64  `json:"room"`
		Organisation *string `json:"organisation"`
		Position     *string `json:"position"`
		IsActive     *bool   `json:"is_active"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	var user data.User

	res := app.DB.First(&user, id)
	if res.RowsAffected == 0 {
		app.notFoundResponse(w, r)
		return
	}

	if res.Error != nil {
		app.serverErrorResponse(w, r, res.Error)
		return
	}

	if input.First_name != nil {
		user.First_name = *input.First_name
	}

	if input.Last_name != nil {
		user.First_name = *input.First_name
	}

	if input.Email != nil {
		user.Email = *input.Email
	}

	if input.PhoneNumber != nil {
		user.PhoneNumber = *input.PhoneNumber
	}

	if input.Password != nil {
		hash, err := bcrypt.GenerateFromPassword([]byte(*input.Password), 12)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		user.Password = string(hash)
	}

	if input.Gender != nil {
		user.Gender = *input.Gender
	}

	if input.RoomID != nil {
		user.RoomID = *input.RoomID
	}

	if input.Organisation != nil {
		user.Organisation = *input.Organisation
	}

	if input.Position != nil {
		user.Position = *input.Position
	}

	if input.IsActive != nil {
		user.IsActive = *input.IsActive
	}

	save := app.DB.Save(&user)

	if save.Error != nil {
		app.serverErrorResponse(w, r, save.Error)
		return
	}

	updated := app.DB.Model(&user).Preload("Role").Preload("Room").First(&user)
	if updated.Error != nil {
		app.serverErrorResponse(w, r, updated.Error)
		return
	}
}

func (app *Application) HandleGetAllUserWithRole(w http.ResponseWriter, r *http.Request) {
	return
}
