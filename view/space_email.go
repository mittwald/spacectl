package view

import (
	"fmt"

	"strings"

	"time"

	"github.com/mittwald/spacectl/client/spaces"
)

type CaughtEmailView struct {
	spaces.CaughtEmail
}

func (v *CaughtEmailView) RenderDate() string {
	return v.Date.Format(time.RFC1123)
}

func (v *CaughtEmailView) RenderSender() string {
	envelopeSender := v.Envelope.MailFrom
	headerSender := v.Header("From")

	if headerSender == "" {
		return fmt.Sprintf("%s (envelope)", envelopeSender)
	}

	if headerSender == envelopeSender {
		return fmt.Sprintf(`%s ("From" header and envelope)`, envelopeSender)
	}

	return fmt.Sprintf(`%s ("From" header), %s (envelope)`, headerSender, envelopeSender)
}

func (v *CaughtEmailView) RenderRecipients(maxCount int) string {
	to := v.Headers("To")
	recipients := make([]string, len(to))

	for i, r := range to {
		recipients[i] = r
	}

	cc := v.Headers("CC")
	for _, r := range cc {
		recipients = append(recipients, r+" (cc)")
	}

	bcc := v.Headers("BCC")
	for _, r := range bcc {
		recipients = append(recipients, r+" (bcc)")
	}

	if len(recipients) > maxCount {
		displayed := recipients[maxCount:]
		rest := len(recipients) - maxCount

		return strings.Join(displayed, ", ") + fmt.Sprintf(" +%d", rest)
	}

	return strings.Join(recipients, ", ")
}
