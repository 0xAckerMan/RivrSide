package data

import "gorm.io/gorm"

type Role struct{
    gorm.Model
    Role string `json:"role_name"`
}
