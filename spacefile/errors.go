package spacefile

import "fmt"

type SyntaxError struct {
	File string
	Inner error
}

func (s SyntaxError) Error() string {
	return fmt.Sprintf(`The file %s could not be parsed due to a syntax error:
    %s

Please have another look at your Spacefile and make sure that the syntax is correct.`, s.File, s.Inner)
}