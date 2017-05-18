package lowlevel

import "fmt"

type ErrUnexpectedStatusCode struct {
	StatusCode int
	Message string
}

func (e ErrUnexpectedStatusCode) Error() string {
	return fmt.Sprintf("unexpected status code %d: %s", e.StatusCode, e.Message)
}