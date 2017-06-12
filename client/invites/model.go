package invites

import (
	"github.com/mittwald/spacectl/client/lowlevel"
	"time"
)

type UserRef struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type TeamRef struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Invite struct {
	Links   lowlevel.LinkList `json:"_links"`
	Actions lowlevel.LinkList `json:"_actions"`

	ID       string    `json:"id"`
	IssuedAt time.Time `json:"issuedAt"`
	Message  string    `json:"message"`
	State    string    `json:"state"`
	Inviter  *UserRef  `json:"inviter"`
	Invitee  *UserRef  `json:"invitee"`
	Team     *TeamRef  `json:"team"`
}
