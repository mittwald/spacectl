package spacefile

import "fmt"

type ErrSpacefileNotFound struct {
	Path string
}

func (n ErrSpacefileNotFound) Error() string {
	return fmt.Sprintf("The Spacefile at %s does not exist, yet. Create one using the \"spacectl space init\" command.", n.Path)
}

type SyntaxError struct {
	File  string
	Inner error
}

func (s SyntaxError) Error() string {
	return fmt.Sprintf(`The file %s could not be parsed due to a syntax error:
    %s

Please have another look at your Spacefile and make sure that the syntax is correct.`, s.File, s.Inner)
}
