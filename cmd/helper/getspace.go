package helper

import (
	"github.com/mittwald/spacectl/client/spaces"
	"github.com/spf13/viper"
	"github.com/mittwald/spacectl/spacefile"
	"fmt"
	"github.com/hashicorp/go-multierror"
	"errors"
	"github.com/mittwald/spacectl/client"
)

func GetSpaceFromContext(args []string, spaceFileName string, api client.SpacesClient) (*spaces.Space, error) {
	if len(args) >= 1 {
		teamID := viper.GetString("teamID")
		if teamID != "" {
			return api.Spaces().GetByTeamAndName(teamID, args[0])
		} else {
			return api.Spaces().GetByID(args[0])
		}
	}

	f, err := spacefile.ParseSpacefile(spaceFileName)
	if err == nil {
		spaceDef := f.Spaces[0]

		space, err := api.Spaces().GetByTeamAndName(spaceDef.TeamID, spaceDef.DNSLabel)
		if err != nil {
			return nil, err
		}

		return space, nil
	}

	if _, ok := err.(spacefile.ErrSpacefileNotFound); ok {
		err := multierror.Append(nil,
			fmt.Errorf("No spacefile found at %s", spaceFileName),
			errors.New("Missing team ID (--team, -t or $SPACES_TEAM_ID)"),
		)
		return nil, err
	}

	return nil, err
}