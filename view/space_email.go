package view

import (
	"fmt"

	"strings"

	"time"

	"io"

	"github.com/fatih/color"
	"github.com/gosuri/uitable"
	"github.com/mittwald/spacectl/client/spaces"
)

var emailSubjectStyle = color.New(color.Underline, color.Bold)

type CaughtEmailView struct {
	spaces.CaughtEmail
}

type CaughtEmailSingleView struct {
	spaces.CaughtEmail
	WithHeaders bool
	AsHTML      bool
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

func (v *CaughtEmailSingleView) Render(out io.Writer) {
	if v.WithHeaders {
		table := uitable.New()
		table.MaxColWidth = 120
		table.Wrap = true

		fmt.Fprintf(out, "ENVELOPE\n")

		table.AddRow("  mailFrom", v.Envelope.MailFrom)
		for i := range v.Envelope.RcptTo {
			table.AddRow("  rcptTo", v.Envelope.RcptTo[i])
		}

		fmt.Fprintln(out, table)

		table = uitable.New()
		table.MaxColWidth = 120
		table.Wrap = true

		fmt.Fprintf(out, "\nHEADER\n")
		for i := range v.Mail.Headers {
			table.AddRow("  "+v.Mail.Headers[i].Name+":", v.Mail.Headers[i].Value)
		}

		fmt.Fprintln(out, table)

		if v.AsHTML && v.Mail.HTML != "" {
			fmt.Fprintf(out, "\nMESSAGE BODY (HTML)\n\n")
		} else {
			fmt.Fprintf(out, "\nMESSAGE BODY (PLAIN TEXT)\n\n")
		}
	}

	fmt.Fprintf(color.Output, emailSubjectStyle.Sprintf(v.Mail.Subject)+"\n\n")

	if v.AsHTML && v.Mail.HTML != "" {
		fmt.Fprintf(out, v.Mail.HTML)
	} else {
		fmt.Fprintf(out, v.Mail.Text)
	}
}
