package data

type Room struct{
    CommonFields
    Number string `json:"room_number" gorm:"unique"`
    Capacity int `json:"capacity"`
    Gender string `json:"gender"`
    Status string `json:"room_status"`
}
