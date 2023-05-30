package core

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/lcrownover/hpcadmin-cli/internal/types"
	"github.com/lcrownover/hpcadmin-cli/internal/util"
)

func CLIUserCreate(username string, email string, firstname string, lastname string) (*types.User, error) {
	u, err := NewUser(username, email, firstname, lastname)
	if err != nil {
		return nil, err
	}
	user, err := CreateUser(u)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%+v\n", user)
	return user, nil
}

func NewUser(username string, email string, firstname string, lastname string) (*types.UserCreate, error) {
	// Validate Username
	if !util.IsAlphanumeric(username) {
		return nil, errors.New("invalid username")
	}
	// Validate Email
	if !util.EmailIsValid(email) {
		return nil, errors.New("invalid email")
	}
	return &types.UserCreate{
		Username:  username,
		Email:     email,
		FirstName: firstname,
		LastName:  lastname,
	}, nil
}

func CreateUser(u *types.UserCreate) (*types.User, error) {
	id := rand.Intn(100000)
	return &types.User{
		Id:        id,
		Username:  u.Username,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}
