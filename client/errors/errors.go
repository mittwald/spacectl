package errors

import "fmt"

type ErrUnauthorized struct {
	Msg string
	Inner error
}

type ErrNested struct {
	Msg string
	Inner error
}

func (a ErrUnauthorized) Error() string {
	return a.Msg
}

func (a ErrNested) Error() string {
	return fmt.Sprintf("%s:\n    %s", a.Msg, a.Inner)
}