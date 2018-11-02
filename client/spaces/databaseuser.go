package spaces

import (
	"fmt"
	"net/url"

	"github.com/mittwald/spacectl/client/errors"
)

// ListDatabaseUserByStage returns an array of DatabaseUserList for the stage
func (c *spacesClient) ListDatabaseUsersByStage(spaceID, stage string) (DatabaseUserList, error) {
	var databaseUserList DatabaseUserList

	listPath := fmt.Sprintf("/spaces/%s/stages/%s/databaseusers", url.PathEscape(spaceID), url.PathEscape(stage))
	err := c.client.Get(listPath, &databaseUserList)
	if err != nil {
		return nil, errors.ErrNested{Inner: err, Msg: fmt.Sprintf("could not access databaseusers for space: %s, stage: %s", spaceID, stage)}
	}

	return databaseUserList, nil
}

// GetDatabaseUser returns information about given databaseUser
func (c *spacesClient) GetDatabaseUser(spaceID, stage, username string) (*DatabaseUser, error) {
	var databaseUser DatabaseUser

	listPath := fmt.Sprintf("/spaces/%s/stages/%s/databaseusers/%s", url.PathEscape(spaceID), url.PathEscape(stage), url.PathEscape(username))
	err := c.client.Get(listPath, &databaseUser)
	if err != nil {
		return nil, errors.ErrNested{Inner: err, Msg: fmt.Sprintf("could not access databaseuser %s for space: %s, stage: %s", username, spaceID, stage)}
	}

	return &databaseUser, nil
}

// CreateDatabaseUser creates a new dbUser for the given stage
func (c *spacesClient) CreateDatabaseUser(spaceID, stage string, dbUser DatabaseUserInput) (*DatabaseUser, error) {
	var newDBUser DatabaseUser
	createPath := fmt.Sprintf("/spaces/%s/stages/%s/databaseusers", url.PathEscape(spaceID), url.PathEscape(stage))
	err := c.client.Post(createPath, dbUser, &newDBUser)
	if err != nil {
		return nil, errors.ErrNested{Inner: err, Msg: fmt.Sprintf("could not create/update dbUser %s", dbUser.UserSuffix)}
	}

	return &newDBUser, err
}

// DeleteDatabaseUser deletes a dbUser from the given stage
// requires the databaseType the user should be deleted from
func (c *spacesClient) DeleteDatabaseUser(spaceID, stage, name, dbType string) error {
	input := DatabaseUserInput{
		Type: dbType,
	}
	out := &Empty{}

	deletePath := fmt.Sprintf("/spaces/%s/stages/%s/databaseusers/%s", url.PathEscape(spaceID), url.PathEscape(stage), url.PathEscape(name))
	err := c.client.DeleteBody(deletePath, input, &out)
	if err != nil {
		return errors.ErrNested{Inner: err, Msg: fmt.Sprintf("could not delete dbUser %s", name)}
	}

	return err
}
