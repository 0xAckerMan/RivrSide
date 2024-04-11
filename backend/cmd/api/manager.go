package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/0xAckerMan/internal/data"
	"golang.org/x/crypto/bcrypt"
)

func (app *Application) HandleCreateNewTenant(w http.ResponseWriter, r *http.Request) {
	var input struct {
		First_name    string `json:"first_name"`
		Last_name     string `json:"last_name"`
		Email         string `json:"email"`
		PhoneNumber   string `json:"phone_number"`
		Password      string `json:"password"`
		Gender        string `json:"gender"`
		RoleID        int64  `json:"role"`
		RoomID        int64  `json:"room"`
		PackageID     int64  `json:"package"`
		Paymentstatus string `json:"payment_status"`
		AmountPaid    int    `json:"amount_paid"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	tenant := data.User{
		First_name:    input.First_name,
		Last_name:     input.Last_name,
		Email:         input.Email,
		PhoneNumber:   input.PhoneNumber,
		RoomID:        input.RoomID,
		PackageID:     input.PackageID,
		Paymentstatus: input.Paymentstatus,
		AmountPaid:    input.AmountPaid,
		IsActive:      true,
	}

	var role data.Role
	roleResult := app.DB.Where(&data.Role{Role: "tenant"}).First(&role)
	if roleResult.RowsAffected == 0 {
		message := "Create a role named tenant first please"
		app.errorResponse(w, r, http.StatusInternalServerError, envelope{"error": message})
		return
	}

	if roleResult.Error != nil {
		app.serverErrorResponse(w, r, roleResult.Error)
		return
	}

	tenant.RoleID = int64(role.ID)

	var packageplan data.PackagePlan

	result := app.DB.First(&packageplan, tenant.PackageID)
	if result.RowsAffected == 0 {
		app.packageNotFoundResponse(w, r)
		return
	}

	if result.Error != nil {
		app.serverErrorResponse(w, r, result.Error)
		return
	}

	var room data.Room

	result = app.DB.First(&room, tenant.RoomID)

	if result.RowsAffected == 0 {
		message := "Invalid room entry, doesnt exist"
		app.errorResponse(w, r, http.StatusNotFound, envelope{"error": message})
		return
	}

	if result.Error != nil {
		app.serverErrorResponse(w, r, result.Error)
		return
	}

	if input.Gender == room.Gender {
		tenant.Gender = input.Gender
	} else {
		message := fmt.Sprintf("%s are not allowed to be assigned in %s rooms", input.Gender, room.Gender)
		app.errorResponse(w, r, http.StatusNotAcceptable, envelope{"error": message})
		return
	}

	tenant.Balance = packageplan.Price - tenant.AmountPaid

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), 12)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	tenant.Password = string(hash)

	res := app.DB.Preload("roles").Preload("package_plans").Preload("rooms").Create(&tenant)
	if res.Error != nil {
		if strings.Contains(res.Error.Error(), "duplicate key value violates unique constraint") {
			app.duplicateRecordResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, res.Error)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"tenant": tenant}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *Application) HandleGetAllTenants(w http.ResponseWriter, r *http.Request) {
    var tenant []data.User

	var role data.Role

	roleResult := app.DB.Where(data.Role{Role: "tenant"}).First(&role)
	if roleResult.RowsAffected == 0 {
		message := "Create a role named tenant first please"
		app.errorResponse(w, r, http.StatusInternalServerError, envelope{"error": message})
		return
	}

	if roleResult.Error != nil {
		app.serverErrorResponse(w, r, roleResult.Error)
		return
	}

	result := app.DB.Where("role_id", role.ID).Preload("Role").Preload("PackagePlan").Preload("Room").Find(&tenant)
    if result.Error!=nil{
        app.serverErrorResponse(w,r,result.Error)
        return
    }

    if result.RowsAffected == 0{
        app.noRecordFoundResponse(w,r)
        return
    }

    err := app.writeJSON(w,http.StatusOK,envelope{"tenants": tenant}, nil)
    if err != nil{
        app.serverErrorResponse(w,r,err)
        return
    }
}

func (app *Application) HandleGetTenantInfo(w http.ResponseWriter, r *http.Request) {
	return
}

func (app *Application) HandleGetAllVaccantRooms(w http.ResponseWriter, r *http.Request) {
	return
}

func (app *Application) HandleGetMaleVaccantRooms(w http.ResponseWriter, r *http.Request) {
	return
}

func (app *Application) HandleGetFemaleVaccantRooms(w http.ResponseWriter, r *http.Request) {
	return
}

func (app *Application) HandleConfirmTenantPayment(w http.ResponseWriter, r *http.Request) {
	return
}

func (app *Application) HandleTransferTenant(w http.ResponseWriter, r *http.Request) {
	return
}

func (app *Application) HandleUpdateRoomStatus(w http.ResponseWriter, r *http.Request) {
	return
}

func (app *Application) HandleGetRoomTenants(w http.ResponseWriter, r *http.Request) {
	return
}
