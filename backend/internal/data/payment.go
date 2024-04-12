package data

import "time"

type Payment struct {
	CommonFields
	UserID      int64        `json:"-"`
	User        *User        `json:"user" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Amount      int          `json:"Amount_paid"`
	Balance     int          `json:"balance"`
	PackageID   int          `json:"-"`
	PackagePlan *PackagePlan `json:"package_plan"  gorm:"foreignKey:PackageID;constraint:OnDelete:CASCADE"`
	Month       time.Month   `json:"month"`
	Year        time.Time    `json:"year"`
}
