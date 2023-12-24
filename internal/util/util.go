package util

import (
	"encoding/json"
	"io"
	"net/mail"
	"regexp"
	"fmt"
	"log/slog"
	"os"

	"github.com/fatih/color"
)

func EmailIsValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func IsAlphanumeric(in string) bool {
	var alphanumeric = regexp.MustCompile("^[a-zA-Z0-9_]*$")
	return alphanumeric.MatchString(in)
}

func DecodeJSON(body io.ReadCloser, v any) error {
	return json.NewDecoder(body).Decode(v)
}

func PrintAndExit(msg string, code int) {
	fmt.Println(msg)
	os.Exit(code)
}

func ConfigureLogging(debug bool) {
	// Set up logging
	lvl := slog.LevelInfo
	if debug {
		lvl = slog.LevelDebug
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: lvl,
	}))
	slog.SetDefault(logger)
}

func SuccessPrint(msg string) {
	color.Green(msg)
}

func InfoPrint(msg string) {
	color.Blue(msg)
}

func ErrorPrint(msg string) {
	color.Red(msg)
}
