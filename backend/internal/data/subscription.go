package data

type Subscription struct {
	CommonFields
    UserID int `json:"-" gorm:"uniqueIndex"`
    User *User `json:"user" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
    PackageID int `json:"-"`
    PackagePlan *PackagePlan `json:"package_plan" gorm:"foreignKey:PackageID;constraint:OnDelete:CASCADE"`
}
