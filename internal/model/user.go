package model

import "time"

type User struct {
	Id int `json:"id"`
	Email string `json:"email"`
	Name string `json:"name"`
	Password string `json:"password,omitempty"`
	Salt string
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
