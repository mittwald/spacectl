package lowlevel

import "fmt"

type ErrUnexpectedStatusCode struct {
	StatusCode int
	Message string
}

func (e ErrUnexpectedStatusCode) Error() string {
	return fmt.Sprintf("unexpected status code %d: %s", e.StatusCode, e.Message)
}

type ErrLinkNotFound struct {
	Rel string
}

func (e ErrLinkNotFound) Error() string {
	return fmt.Sprintf("link \"%s\" does not exist", e.Rel)
}