package teams

import (
	"fmt"
	"github.com/mittwald/spacectl/client/errors"
	"github.com/mittwald/spacectl/client/lowlevel"
)

func (t *teamsClient) Delete(idOrLabel string) error {
	var team Team
	var result lowlevel.Message

	url := fmt.Sprintf("/v1/teams/%s", idOrLabel)
	err := t.client.Get(url, &team)
	if err != nil {
		return fmt.Errorf("could not delete team: %s", err)
	}

	link, err := team.Actions.GetLinkByRel("delete")
	if err != nil {
		switch err.(type) {
		case lowlevel.ErrLinkNotFound:
			return errors.ErrUnauthorized{Msg: "You are not allowed to delete this team.", Inner: err}
		default:
			return err
		}
	}

	return link.Delete(t.client, &result)
}
