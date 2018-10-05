package spaces

import "time"

type DatabaseUser struct {
	User        string    `json:"user"`
	CreatedAt   time.Time `json:"createdAt"`
	ModifiedAt  time.Time `json:"modifiedAt"`
	Status      string    `json:"status"`
	Type        string    `json:"type"`
	Stage       string    `json:"stage"`
}

type DatabaseUserList []DatabaseUser
