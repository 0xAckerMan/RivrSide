package data

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	First_name    string     `json:"first_name"`
	Last_name     string     `json:"last_name"`
	Email         string     `json:"email"`
	PhoneNumber   string     `json:"phone_number"`
	Password      string     `json:"-"`
	Gender        string     `json:"gender"`
	Role          *Role      `json:"role"`
	Room          *Room      `json:"room"`
	Organisation  string     `json:"organisation"`
	Position      string     `json:"position"`
	Package       *Package   `json:"package"`
	Month         time.Month `json:"current_month"`
	Paymentstatus string     `json:"payment_status"`
	Amount        int        `json:"amount"`
	Balance       int        `json:"balance"`
	Isadmin       bool       `json:"-"`
	Ismanager     bool       `json:"-"`
	IsActive      bool       `json:"is_active"`
}
