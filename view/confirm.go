package view

import (
	"github.com/fatih/color"
	"github.com/mittwald/spacectl/view/confirm"
)

func Confirm(msg string, body string) (bool, error) {
	if msg == "" {
		msg = "This action is destructive!"
	}

	c := confirm.Confirmation{
		Color: color.New(color.FgRed),
		Message: msg,
		Title: "Warning",
		Body: body,
		Prompt: confirm.DefaultPrompt,
	}

	return c.DoPrompt(color.Output)
}
