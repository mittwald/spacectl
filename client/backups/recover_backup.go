package backups

import "github.com/mittwald/spacectl/client/errors"

type recoveryRequest struct {
	Files interface{} `json:"files"`
	Databases interface{} `json:"databases"`
	Metadata interface{} `json:"metadata"`
}

func (r *RecoverySpec) buildRequest() interface{} {
	switch r.Type {
	case RecoverAll:
		return "all"
	case RecoverNone:
		return "none"
	case RecoverSpecific:
		return r.Items
	}

	return "all"
}

func (c *backupClient) Recover(backupID string, files RecoverySpec, databases RecoverySpec, metadata RecoverySpec) (*Recovery, error) {
	backup, err := c.Get(backupID)
	if err != nil {
		return nil, err
	}

	recoverLink, err := backup.Actions.GetLinkByRel("recover")
	if err != nil {
		return nil, errors.ErrUnauthorized{Inner: err, Msg: "recovery is not available for this backup"}
	}

	res := Recovery{}
	req := recoveryRequest{
		Files: files.buildRequest(),
		Databases: databases.buildRequest(),
		Metadata: metadata.buildRequest(),
	}

	err = recoverLink.Execute(c.client, req, res)
	if err != nil {
		return nil, errors.ErrNested{Inner: err, Msg: "could not start backup recovery"}
	}

	return &res, nil
}