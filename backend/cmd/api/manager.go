package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/0xAckerMan/internal/data"
	"golang.org/x/crypto/bcrypt"
)

func (app *Application) HandleCreateNewTenant(w http.ResponseWriter, r *http.Request) {
	var input data.CreateTenant

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	tenant := data.User{
		First_name:  input.First_name,
		Last_name:   input.Last_name,
		Email:       input.Email,
		PhoneNumber: input.PhoneNumber,
		RoomID:      input.RoomID,
		IsActive:    true,
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

	var room data.Room

	result := app.DB.First(&room, tenant.RoomID)

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

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), 12)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	tenant.Password = string(hash)

	var occupants []data.User

	capacity := app.DB.Where("room_id = ?", room.ID).Find(&occupants)
	if capacity.Error != nil {
		app.serverErrorResponse(w, r, capacity.Error)
		return
	}

	if capacity.RowsAffected >= int64(room.Capacity) {
		err := app.writeJSON(w, http.StatusForbidden, envelope{"response": envelope{"error": "Ooops! Sorry the room capacity is full at the moment"}}, nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		return
	}

	res := app.DB.Create(&tenant)
	if res.Error != nil {
		if strings.Contains(res.Error.Error(), "duplicate key value violates unique constraint") {
			app.duplicateRecordResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, res.Error)
		return
	}

	if err := app.DB.Model(&tenant).Preload("Role").Preload("Room").First(&tenant).Error; err != nil {
		app.serverErrorResponse(w, r, err)
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

	result := app.DB.Where("role_id", role.ID).Preload("Role").Preload("Room").Find(&tenant)
	if result.Error != nil {
		app.serverErrorResponse(w, r, result.Error)
		return
	}

	if result.RowsAffected == 0 {
		app.noRecordFoundResponse(w, r)
		return
	}

	err := app.writeJSON(w, http.StatusOK, envelope{"tenants": tenant}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *Application) HandleGetTenantInfo(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDparam(w, r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var tenant data.User

	result := app.DB.Preload("Role").Preload("Room").First(&tenant, id)
	if result.RowsAffected == 0 {
		app.noRecordFoundResponse(w, r)
		return
	}

	if result.Error != nil {
		app.serverErrorResponse(w, r, result.Error)
		return
	}

	tenantInfo := data.TenantInfo{
        User: tenant,  
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"tenant": tenantInfo}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *Application) HandleGetAllActiveTenants(w http.ResponseWriter, r *http.Request) {
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

	result := app.DB.Where("role_id = ? AND is_active = ?", role.ID, "true").Preload("Role").Preload("Room").Find(&tenant)
	if result.Error != nil {
		app.serverErrorResponse(w, r, result.Error)
		return
	}

	if result.RowsAffected == 0 {
		app.noRecordFoundResponse(w, r)
		return
	}

	err := app.writeJSON(w, http.StatusOK, envelope{"tenants": tenant}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *Application) HandleGetAllVacantRooms(w http.ResponseWriter, r *http.Request) {
    type VacantRoomInfo struct {
        data.Room
        Vacancies int `json:"vacancies"`
    }

    // Query rooms and calculate vacancies per room
    var vacantRoomInfos []VacantRoomInfo
    result := app.DB.
        Table("rooms").
        Select("rooms.*, (rooms.capacity - COUNT(users.id)) AS vacancies").
        Joins("LEFT JOIN users ON rooms.id = users.room_id").
        Group("rooms.id").
        Having("(rooms.capacity - COUNT(users.id)) > 0").
        Find(&vacantRoomInfos)
    if result.Error != nil {
        app.serverErrorResponse(w, r, result.Error)
        return
    }

    err := app.writeJSON(w, http.StatusOK, envelope{"vacant_rooms": vacantRoomInfos}, nil)
    if err != nil {
        app.serverErrorResponse(w, r, err)
        return
    }
}

func (app *Application) HandleGetMaleVacantRooms(w http.ResponseWriter, r *http.Request) {
    type VacantRoomInfo struct {
        data.Room
        Vacancies int `json:"vacancies"`
    }

    // Query rooms and calculate vacancies per room for male occupants
    var vacantRoomInfos []VacantRoomInfo
    result := app.DB.
        Table("rooms").
        Select("rooms.*, (rooms.capacity - COUNT(users.id)) AS vacancies").
        Joins("LEFT JOIN users ON rooms.id = users.room_id").
        Where("users.gender = ?", "male").
        Group("rooms.id").
        Having("(rooms.capacity - COUNT(users.id)) > 0").
        Find(&vacantRoomInfos)
    if result.Error != nil {
        app.serverErrorResponse(w, r, result.Error)
        return
    }

    err := app.writeJSON(w, http.StatusOK, envelope{"male_vacant_rooms": vacantRoomInfos}, nil)
    if err != nil {
        app.serverErrorResponse(w, r, err)
        return
    }
}

func (app *Application) HandleGetFemaleVacantRooms(w http.ResponseWriter, r *http.Request) {
    type VacantRoomInfo struct {
        data.Room
        Vacancies int `json:"vacancies"`
    }

    // Query rooms and calculate vacancies per room for male occupants
    var vacantRoomInfos []VacantRoomInfo
    result := app.DB.
        Table("rooms").
        Select("rooms.*, (rooms.capacity - COUNT(users.id)) AS vacancies").
        Joins("LEFT JOIN users ON rooms.id = users.room_id").
        Where("users.gender = ?", "female").
        Group("rooms.id").
        Having("(rooms.capacity - COUNT(users.id)) > 0").
        Find(&vacantRoomInfos)
    if result.Error != nil {
        app.serverErrorResponse(w, r, result.Error)
        return
    }

    err := app.writeJSON(w, http.StatusOK, envelope{"female_vacant_rooms": vacantRoomInfos}, nil)
    if err != nil {
        app.serverErrorResponse(w, r, err)
        return
    }
}


func (app *Application) HandleConfirmTenantPayment(w http.ResponseWriter, r *http.Request) {
	return
}

func (app *Application) HandleTransferTenant(w http.ResponseWriter, r *http.Request) {
	return
}

func (app *Application) HandleUpdateRoomStatus(w http.ResponseWriter, r *http.Request) {
    id, err := app.readIDparam(w,r)
    if err != nil{
        app.serverErrorResponse(w,r,err)
        return
    }
    var input struct{
        Status *string `json:"room_status"`
    }
	
    err = app.readJSON(w,r,&input)
    if err != nil{
        app.serverErrorResponse(w,r,err)
        return
    }

    var room data.Room

    res := app.DB.First(&room, id)
    if res.RowsAffected == 0{
        app.noRecordFoundResponse(w,r)
        return
    }

    if res.Error != nil{
        app.serverErrorResponse(w,r,res.Error)
        return
    }

    if input.Status != nil{
        room.Status = *input.Status
    }

    result := app.DB.Save(&room)
    if result.Error != nil{
        app.serverErrorResponse(w,r,result.Error)
        return
    }

    err = app.writeJSON(w,http.StatusOK,envelope{"room": room}, nil)
    if err != nil{
        app.serverErrorResponse(w,r,result.Error)
        return
    }
    
}

func (app *Application) HandleGetRoomTenants(w http.ResponseWriter, r *http.Request) {
    id, err := app.readIDparam(w,r)
    if err != nil{
        app.serverErrorResponse(w,r,err)
        return
    }

    var room data.Room
    res:= app.DB.First(&room, id)
    if res.RowsAffected == 0{
        message:= "Ooops that room doesnt exist"
        app.errorResponse(w,r,http.StatusNotFound,envelope{"error":message})
        return
    }

    if res.Error != nil{
        app.serverErrorResponse(w,r,res.Error)
        return
    }

    var tenants []data.User

    results := app.DB.Preload("Role").Preload("Room").Where(data.User{RoomID: int64(room.ID)}).Find(&tenants)
    if results.RowsAffected == 0 {
        app.noRecordFoundResponse(w,r)
        return
    }

    if results.Error != nil{
        app.serverErrorResponse(w,r,results.Error)
        return
    }

    err = app.writeJSON(w,http.StatusOK,envelope{"tenants": tenants}, nil)
    if err != nil{
        app.serverErrorResponse(w,r,err)
        return
    }
}
