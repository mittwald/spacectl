package invites

type Invite struct {
	ID      string `json:"id"`
	Message string `json:"message"`
	State   string `json:"state"`
}
