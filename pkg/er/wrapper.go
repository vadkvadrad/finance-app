package er

import (
	"fmt"
)

const (
	ErrNotAuthorized        = "user are not authorized"
	ErrWrongUserCredentials = "wrong user credentials"
	ErrUserExists           = "user already exists"
	ErrUserNotExists        = "user not exists"
	ErrUserNotVerified    = "user not verified"
)

func Wrap(msg string, err error) error {
	return fmt.Errorf("%s: %s", msg, err)
}
