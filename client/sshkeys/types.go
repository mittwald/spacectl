package sshkeys

import "time"

type sshKeyInput struct {
	CipherAlgorithm string    `json:"cipherAlgorithm"`
	Comment         string    `json:"comment,omitempty"`
	Key             []byte    `json:"key"`
}

type SSHKey struct {
	ID              string    `json:"id,omitempty"`
	CreatedAt       time.Time `json:"createdAt,omitempty"`
	CipherAlgorithm string    `json:"cipherAlgorithm"`
	Comment         string    `json:"comment,omitempty"`
	Key             []byte    `json:"key"`
}
