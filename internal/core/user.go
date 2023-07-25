package core

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/lcrownover/hpcadmin-cli/internal/types"
	"github.com/lcrownover/hpcadmin-cli/internal/util"
)

func CLIUserCreate(username string, email string, firstname string, lastname string) (*types.UserResponse, error) {
	u, err := NewUserRequest(username, email, firstname, lastname)
	if err != nil {
		return nil, err
	}
	user, err := CreateUser(u)
	if err != nil {
        fmt.Println(err)
		return nil, err
	}
	return user, nil
}

func CLIUserShow(username, string) (*types.UserResponse, error) {
    u, err := GetUser(username)
    if err != nil {
        return nil, err
    }
    return nil, nil
}

func GetUser(username string) (*types.UserResponse, error) {
    // TODO(lcrown): get the user by username
    // should it take a username, or a "queryable" interface?
    // i assume we should be able to get users by username, email, or ID.
    return nil, nil
}

// NewUserRequest validates user input and returns a UserRequest struct
func NewUserRequest(username string, email string, firstname string, lastname string) (*types.UserRequest, error) {
	// Validate Username
	if !util.IsAlphanumeric(username) {
		return nil, errors.New("invalid username")
	}
	// Validate Email
	if !util.EmailIsValid(email) {
		return nil, errors.New("invalid email")
	}
	return &types.UserRequest{
		Username:  username,
		Email:     email,
		FirstName: firstname,
		LastName:  lastname,
	}, nil
}

func CreateUser(u *types.UserRequest) (*types.UserResponse, error) {

    // TODO(lcrown): request new user from API
    // This needs to be completely refactored, but lets get it working
    req, err := NewAPIRequest("users", "POST", u)
    if err != nil {
        return nil, err
    }
    baseUrl := "http://localhost:3333"
    url := fmt.Sprintf("%s/%s", baseUrl, req.Endpoint)
    // fmt.Printf("body: %v\n", string(req.Body))
    body := bytes.NewReader(req.Body)
    newReq, err := http.NewRequest(req.Method, url, body)
    if err != nil {
        return nil, err
    }
    newReq.Header.Set("Content-Type", "application/json")
    client := &http.Client{}
    resp, err := client.Do(newReq)
    if err != nil {
        return nil, err
    }
    fmt.Println(resp.StatusCode)
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusCreated {
        text, _ := io.ReadAll(resp.Body)
        return nil, errors.New(string(text))
    }
    var user types.UserResponse

    err = util.DecodeJSON(resp.Body, &user)
    if err != nil {
        return nil, err
    }
    return &user, nil
}
