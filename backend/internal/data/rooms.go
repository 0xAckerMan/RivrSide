package data

import "gorm.io/gorm"

type Room struct{
    gorm.Model
    Number string `json:"room_number"`
    Capacity int `json:"capacity"`
    Gender string `json:"gender"`
    Status string `json:"room_status"`
}
