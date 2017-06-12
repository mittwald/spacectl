package lowlevel

import "fmt"

type Message struct {
	Message string `json:"msg"`
	Error   string `json:"err"`
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

type Link struct {
	Href string `json:"href"`
	Rel string `json:"rel"`
	Method string `json:"method"`
}

func (l Link) Get(client *SpacesLowlevelClient, result interface{}) error {
	return client.Get(l.Href, result)
}

func (l Link) Post(client *SpacesLowlevelClient, body interface{}, result interface{}) error {
	return client.Post(l.Href, body, result)
}

type LinkList []Link

func (l LinkList) HasLink(rel string) bool {
	for i := range l {
		if l[i].Rel == rel {
			return true
		}
	}

	return false
}

func (l LinkList) GetLinkByRel(rel string) (*Link, error) {
	for i := range l {
		if l[i].Rel == rel {
			return &l[i], nil
		}
	}

	return nil, ErrLinkNotFound{rel}
}
