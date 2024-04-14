package main

import (
	"net/http"
	"strings"

	"github.com/0xAckerMan/internal/data"
)

func (app *Application) HandleCreateRoom(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Number   string `json:"room_number"`
		Capacity int    `json:"capacity"`
		Gender   string `json:"gender"`
		Status   string `json:"room_status"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	room := data.Room{
		Number:   input.Number,
		Capacity: input.Capacity,
		Gender:   input.Gender,
		Status:   input.Status,
	}

	err = app.DB.Create(&room).Error
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			app.duplicateRecordResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"room": room}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *Application) HandleGetSingleRoom(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDparam(w, r)
	if err != nil {
		app.noRecordFoundResponse(w, r)
		return
	}

	type Room struct {
		data.Room
		Vacancies int `json:"vacancies"`
	}

	var room Room

	result := app.DB.Table("rooms").
		Select("rooms.*, rooms.capacity - COUNT(users.id) AS vacancies").
		Joins("LEFT JOIN users ON rooms.id = users.room_id").
		Where("rooms.id = ?", id).
		Where("rooms.deleted_at IS NULL").
		Group("rooms.id").Having("(rooms.capacity - COUNT(users.id)) >= 0").
		Order("rooms.id").
		First(&room)

	if result.RowsAffected == 0 {
		app.noRecordFoundResponse(w, r)
		return
	}

	if result.Error != nil {
		app.serverErrorResponse(w, r, result.Error)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"room": room}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *Application) HandleGetAllRooms(w http.ResponseWriter, r *http.Request) {
	type RoomInfo struct {
		data.Room
		Vacancies int `json:"vacancies"`
	}

	// Query rooms and calculate vacancies per room
	var RoomInfos []RoomInfo
	result := app.DB.
		Table("rooms").
		Select("rooms.*, (rooms.capacity - COUNT(users.id)) AS vacancies").
		Joins("LEFT JOIN users ON rooms.id = users.room_id").
		Group("rooms.id").
		Find(&RoomInfos)
	if result.Error != nil {
		app.serverErrorResponse(w, r, result.Error)
		return
	}

	err := app.writeJSON(w, http.StatusOK, envelope{"rooms": RoomInfos}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *Application) HandlePartialUpdateRoom(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDparam(w, r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var room data.Room

	result := app.DB.First(&room, id)
	if result.Error != nil {
		app.notFoundResponse(w, r)
		return
	}

	var input struct {
		Number   *string `json:"room_number"`
		Capacity *int    `json:"capacity"`
		Gender   *string `json:"gender"`
		Status   *string `json:"room_status"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if input.Number != nil {
		room.Number = *input.Number
	}

	if input.Capacity != nil {
		room.Capacity = *input.Capacity
	}

	if input.Gender != nil {
		room.Gender = *input.Gender
	}

	if input.Status != nil {
		room.Status = *input.Status
	}

	err = app.DB.Save(&room).Error
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"room": room}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *Application) HandleDeleteRoom(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDparam(w, r)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	var room data.Room

	result := app.DB.Delete(&room, id)
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
