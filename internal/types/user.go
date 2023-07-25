package types

import (
	"encoding/json"
	"time"
)

type UserRequest struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

func (u *UserRequest) ToBytes() ([]byte, error) {
    return json.Marshal(u)
}

type UserResponse struct {
	Id        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FirstName string    `json:"firstname"`
	LastName  string    `json:"lastname"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserStub struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
}
