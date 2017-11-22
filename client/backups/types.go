package backups

import (
	"github.com/mittwald/spacectl/client/lowlevel"
	"time"
)

type Backup struct {
	Links   lowlevel.LinkList `json:"_links"`
	Actions lowlevel.LinkList `json:"_actions"`

	ID          string    `json:"id"`
	StartedAt   time.Time `json:"startedAt"`
	CompletedAt time.Time `json:"completedAt"`
	Status      string    `json:"status"`
	Keep        bool      `json:"keep"`
	Description string    `json:"description"`

	Stage *StageRef `json:"stage"`
	Space *SpaceRef `json:"space"`
}

type BackupRef struct {
	ID   string `json:"id"`
	HREF string `json:"href"`
}

type StageRef struct {
	Name string `json:"name"`
	HREF string `json:"href"`
}

type SpaceRef struct {
	ID   string `json:"id"`
	HREF string `json:"href"`
}

type Recovery struct {
	Links   lowlevel.LinkList `json:"_links"`
	Actions lowlevel.LinkList `json:"_actions"`

	ID          string    `json:"id"`
	StartedAt   time.Time `json:"startedAt"`
	CompletedAt time.Time `json:"completedAt"`
	Status      string    `json:"status"`

	Backup *BackupRef `json:"backup,omitempty"`
	Stage  *StageRef  `json:"stage,omitempty"`
}
