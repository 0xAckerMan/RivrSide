package data

import "gorm.io/gorm"

type Payment struct {
	gorm.Model
	User    *User    `json:"user"`
	Amount  int      `json:"Amount"`
	Package *Package `json:"package"`
}
