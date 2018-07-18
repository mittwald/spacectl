package spaces

import (
	"strings"
	"time"
)

type Envelope struct {
	MailFrom string   `json:"mailFrom"`
	RcptTo   []string `json:"rcptTo"`
}

type Header struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Mail struct {
	Subject string   `json:"subject"`
	Headers []Header `json:"headers"`
	Text    string   `json:"text"`
	HTML    string   `json:"html"`
	To      string   `json:"to"`
	From    string   `json:"from"`
	CC      string   `json:"cc"`
	BCC     string   `json:"bcc"`
}

type CaughtEmail struct {
	ID       string    `json:"id"`
	Envelope Envelope  `json:"envelope"`
	Date     time.Time `json:"date"`
	Mail     Mail      `json:"mail"`
}

func (c *CaughtEmail) Header(name string) string {
	name = strings.ToLower(name)
	for i := range c.Mail.Headers {
		if strings.ToLower(c.Mail.Headers[i].Name) == name {
			return c.Mail.Headers[i].Value
		}
	}

	return ""
}

func (c *CaughtEmail) Headers(name string) []string {
	name = strings.ToLower(name)
	values := make([]string, 0)

	for i := range c.Mail.Headers {
		if strings.ToLower(c.Mail.Headers[i].Name) == name {
			values = append(values, c.Mail.Headers[i].Value)
		}
	}

	return values
}

type CaughtEmailList []CaughtEmail

func (l CaughtEmailList) ByID(id string) *CaughtEmail {
	for i := range l {
		if l[i].ID == id {
			return &l[i]
		}
	}

	return nil
}
