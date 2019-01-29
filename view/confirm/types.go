package confirm

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type Confirmation struct {
	Prompt  string
	Title   string
	Message string
	Body    interface{}
	Color   *color.Color
}

func (c *Confirmation) title() string {
	if c.Title == "" {
		return strings.ToUpper(DefaultTitle)
	}

	return strings.ToUpper(c.Title)
}

func (c *Confirmation) DoPrompt(out io.Writer) (bool, error) {
	msg := c.Message
	col := c.Color

	if msg == "" {
		msg = "This action is destructive!"
	}

	if col == nil {
		col = color.New(color.FgRed)
	}

	if c.Prompt != "" {
		fmt.Fprintln(color.Output, c.Prompt)
	}

	msgLines := strings.Split(msg, "\n")

	col.Fprintln(out, c.title())

	for _, l := range msgLines {
		col.Fprintln(out, "  "+l)
	}

	col.Fprintln(out, "  Are you ABSOLUTELY sure you want to continue?")

	fmt.Println("")

	switch b := c.Body.(type) {
	case string:
		fmt.Println(b)
		fmt.Println("")
	case io.Reader:
		s, err := ioutil.ReadAll(b)
		if err != nil {
			return false, err
		}

		fmt.Println(strings.TrimSpace(string(s)))
		fmt.Println("")
	case nil:
		break
	default:
		return false, fmt.Errorf("unknown body type: %T", b)
	}

	for {
		fmt.Fprint(color.Output, "Enter "+color.YellowString("y")+" or "+color.YellowString("n")+": ")

		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')

		input = strings.TrimRight(input, "\n")

		if input == "y" {
			return true, nil
		} else if input == "n" {
			return false, nil
		}
	}
}
