package backups

import (
	"github.com/mittwald/spacectl/client/lowlevel"
	"log"
	"time"
)

type RecoverySpecType int

const (
	RecoverAll RecoverySpecType = iota
	RecoverNone
	RecoverSpecific
)

type RecoverySpec struct {
	Type  RecoverySpecType
	Items []string
}

type BackupListOptions struct {
	OnlyKeep bool
	Since    time.Time
}

type BackupClient interface {
	ListForSpace(spaceID string, opts *BackupListOptions) ([]Backup, error)
	ListForStage(spaceID, stage string, opts *BackupListOptions) ([]Backup, error)
	Get(backupID string) (*Backup, error)
	Create(spaceID string, stage string, keep bool, description string) (*Backup, error)
	Delete(backupID string) error
	Recover(backupID string, stage string, files RecoverySpec, databases RecoverySpec, metadata RecoverySpec) (*Recovery, error)
}

type RecoveryClient interface {
	ListForSpace(spaceID string) ([]Recovery, error)
	ListForBackupID(backupID string) ([]Recovery, error)
	ListForBackup(backup *Backup) ([]Recovery, error)
}

func NewBackupClient(c *lowlevel.SpacesLowlevelClient, l *log.Logger) BackupClient {
	return &backupClient{c, l}
}

func NewRecoveryClient(c *lowlevel.SpacesLowlevelClient, l *log.Logger) RecoveryClient {
	return &recoveryClient{c, l}
}

type backupClient struct {
	client *lowlevel.SpacesLowlevelClient
	logger *log.Logger
}

type recoveryClient struct {
	client *lowlevel.SpacesLowlevelClient
	logger *log.Logger
}
