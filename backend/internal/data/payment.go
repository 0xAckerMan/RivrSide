package data

import "time"

type Payment struct {
	CommonFields
	UserID         int64         `json:"-"`
	User           *User         `json:"user" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Amount         int           `json:"Amount_paid"`
	Balance        int           `json:"balance"`
	SubscriptionID int           `json:"-"`
	Subscription   *Subscription `json:"package_plan"  gorm:"foreignKey:SubscriptionID;constraint:OnDelete:CASCADE"`
	Month          time.Month    `json:"month"`
	Year           int     `json:"year"`
}

type MakePayment struct {
	Amount int `json:"amount"`
}
