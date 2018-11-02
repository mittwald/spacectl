package spaces

import "time"

type DatabaseUser struct {
	User       string    `json:"user"`
	CreatedAt  time.Time `json:"createdAt"`
	ModifiedAt time.Time `json:"modifiedAt"`
	Status     string    `json:"status"`
	Type       string    `json:"type"`
	Stage      StageRef  `json:"stage"`
}

type DatabaseUserList []DatabaseUser

type DatabaseUserInput struct {
	UserSuffix string `json:"userSuffix"`
	Password   string `json:"password"`
	External   string `json:"external"`
	Type       string `json:"type"`
}
