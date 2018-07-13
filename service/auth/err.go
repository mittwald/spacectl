package auth

import "fmt"

type AuthErr struct {
	inner error
}

func (e AuthErr) Error() string {
	return fmt.Sprintf("authentication error: %s", e.inner.Error())
}

type InvalidCredentialsErr struct{}

func (e InvalidCredentialsErr) Error() string {
	return "invalid credentials"
}
