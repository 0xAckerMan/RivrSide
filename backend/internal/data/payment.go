package data

type Payment struct {
    CommonFields
    UserID int64 `json:"-"`
	User    *User    `json:"user" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Amount  int      `json:"Amount"`
    PackageID int `json:"-"`
	PackagePlan *PackagePlan `json:"package_plan"  gorm:"foreignKey:PackageID;constraint:OnDelete:CASCADE"`
}
