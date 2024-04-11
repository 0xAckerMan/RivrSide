package main

import (
	"net/http"

	"github.com/0xAckerMan/internal/data"
)

func (app *Application) HandleCreateRole(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Role string `json:"role_name"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	role := data.Role{
		Role: input.Role,
	}

	err = app.DB.Create(&role).Error
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"role": role}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *Application) HandleGetAllRole(w http.ResponseWriter, r *http.Request) {
	var roles []data.Role

	result := app.DB.Find(&roles)
	if result.RowsAffected == 0 {
		app.noRecordFoundResponse(w, r)
		return
	}

	if result.Error != nil {
		app.serverErrorResponse(w, r, result.Error)
		return
	}

	err := app.writeJSON(w, http.StatusOK, envelope{"roles": roles}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *Application) HandleUpdateRole(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDparam(w, r)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	var role data.Role

	err = app.DB.First(&role, id).Error
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	var input struct {
		Role *string `json:"role_name"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if input.Role != nil {
		role.Role = *input.Role
	}

	err = app.DB.Save(&role).Error
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"role": role}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *Application) HandleDeleteRole(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDparam(w, r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var role data.Role

	result := app.DB.Delete(&role, id)
	if result.RowsAffected == 0 {
		app.noRecordFoundResponse(w, r)
		return
	}

	if result.Error != nil {
		app.serverErrorResponse(w, r, result.Error)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"success": "record deleted successfuly"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

