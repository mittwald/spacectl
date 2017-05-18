package lowlevel

import "fmt"

type Message struct {
	Message string `json:"msg"`
	Error string `json:"err"`
}

func (m Message) String() string {
	if m.Message != "" && m.Error != "" {
		return fmt.Sprintf("%s (%s)", m.Message, m.Error)
	}

	if m.Message != "" {
		return m.Message
	}

	if m.Error != "" {
		return m.Error
	}

	return "Unknown"
}