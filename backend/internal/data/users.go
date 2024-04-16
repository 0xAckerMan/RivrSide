package data

import (
	"time"
)

type User struct {
	CommonFields
	First_name   string `json:"first_name" gorm:"not null"`
	Last_name    string `json:"last_name" gorm:"not null"`
	Email        string `json:"email" gorm:"uniqueIndex;not null"`
	PhoneNumber  string `json:"phone_number" gorm:"uniqueIndex;not null"`
	Password     string `json:"-" gorm:"not null"`
	Gender       string `json:"gender"`
	RoleID       int64  `json:"-"`
	Role         *Role  `json:"role" gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE"`
    RoomID       int64  `json:"-" gorm:"default:null"`
	Room         *Room  `json:"room" gorm:"foreignKey:RoomID;constraint:OnDelete:CASCADE"`
	Organisation string `json:"organisation"`
	Position     string `json:"position"`
	Isadmin      bool   `json:"-" gorm:"default:false"`
	Ismanager    bool   `json:"-" gorm:"default:false"`
	IsActive     bool   `json:"is_active"`
}

type CreateTenant struct {
	First_name  string `json:"first_name"`
	Last_name   string `json:"last_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
	Gender      string `json:"gender"`
	RoleID      int64  `json:"role"`
	RoomID      int64  `json:"room"`
}

type UpdateTenant struct {
	First_name  *string `json:"first_name"`
	Last_name   *string `json:"last_name"`
	Email       *string `json:"email"`
	PhoneNumber *string `json:"phone_number"`
	Password    *string `json:"password"`
	Gender      *string `json:"gender"`
	RoleID      *int64  `json:"role"`
	RoomID      *int64  `json:"room"`
}

type TenantInfo struct {
	User
	PackageID     int64        `json:"package_id"`
	PackagePlan   *PackagePlan `json:"package_plan"`
	Month         time.Month   `json:"current_month"`
	Paymentstatus string       `json:"payment_status"`
	AmountPaid    int          `json:"amount_paid"`
	Balance       int          `json:"balance"`
}

type Manager struct {
	CommonFields
	First_name   string `json:"first_name" gorm:"not null"`
	Last_name    string `json:"last_name" gorm:"not null"`
	Email        string `json:"email" gorm:"uniqueIndex;not null"`
	PhoneNumber  string `json:"phone_number" gorm:"uniqueIndex;not null"`
	Password     string `json:"-" gorm:"not null"`
	Gender       string `json:"gender"`
	RoleID       int64  `json:"-"`
	Role         *Role  `json:"role" gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE"`
	Organisation string `json:"organisation"`
	Position     string `json:"position"`
	Ismanager    bool   `json:"-" gorm:"default:false"`
	IsActive     bool   `json:"is_active"`
}

type CreateManager struct {
	First_name   string `json:"first_name"`
	Last_name    string `json:"last_name"`
	Email        string `json:"email"`
	PhoneNumber  string `json:"phone_number"`
	Password     string `json:"password"`
	Gender       string `json:"gender"`
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
