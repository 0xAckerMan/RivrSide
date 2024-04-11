package main

import (
	"net/http"

	"github.com/0xAckerMan/internal/data"
)

func (app *Application) HandleCreateNewTenant(w http.ResponseWriter, r *http.Request) {
	var input struct {
		First_name    string       `json:"first_name"`
		Last_name     string       `json:"last_name"`
		Email         string       `json:"email"`
		PhoneNumber   string       `json:"phone_number"`
		Password      string       `json:"password"`
		Gender        string       `json:"gender"`
		RoleID        int64        `json:"role"`
		RoomID        int64        `json:"room"`
		PackageID     int64        `json:"package"`
		Paymentstatus string       `json:"payment_status"`
		AmountPaid    int          `json:"amount_paid"`
	}
    
    err := app.readJSON(w,r,&input)
    if err != nil {
        app.serverErrorResponse(w,r,err)
        return
    }

    tenant := data.User{
        First_name: input.First_name,
        Last_name: input.Last_name,
        Email: input.Email,
        PhoneNumber: input.PhoneNumber,
        Gender: input.Gender,
        RoleID: input.RoleID,
        RoomID: input.RoomID,
        PackageID: input.PackageID,
        Paymentstatus: input.Paymentstatus,
        AmountPaid: input.AmountPaid,
        IsActive: true,
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
