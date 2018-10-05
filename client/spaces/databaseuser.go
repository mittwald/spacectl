package spaces

import (
	"fmt"
	"github.com/mittwald/spacectl/client/errors"
	"net/url"
)

// ListDatabaseUserByStage returns an array of DatabaseUserList for the stage
func (c *spacesClient) ListDatabaseUsersByStage(spaceID, stage string) (DatabaseUserList, error) {
	var databaseUserList DatabaseUserList

	listPath := fmt.Sprintf("/spaces/%s/stages/%s/databaseusers", url.PathEscape(spaceID), url.PathEscape(stage))
	err := c.client.Get(listPath, &listPath)
	if err != nil {
		return nil, errors.ErrNested{Inner: err, Msg: fmt.Sprintf("could not access databaseusers for space: %s, stage: %s", spaceID, stage)}
	}

	return databaseUserList, nil
}
