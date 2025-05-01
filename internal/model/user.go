package model

import "time"

type User struct {
	Id int `json:"id"`
	Email string `json:"email" validate:"required"`
	Name string `json:"name" validate:"required,min=5,max=20"`
	Password string `json:"password,omitempty" validate:"required,min=8"`
	Phone string `json:"phone" validate:"required,min=8"`
	Salt string
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}