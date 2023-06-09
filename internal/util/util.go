package util

import (
	"net/mail"
	"regexp"
)

func EmailIsValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func IsAlphanumeric(in string) bool {
	var alphanumeric = regexp.MustCompile("^[a-zA-Z0-9_]*$")
	return alphanumeric.MatchString(in)
}
