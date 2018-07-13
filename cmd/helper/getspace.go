package helper

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/mittwald/spacectl/client"
	"github.com/mittwald/spacectl/client/spaces"
	"github.com/mittwald/spacectl/spacefile"
	"github.com/spf13/viper"
)

func GetSpaceFromContext(args []string, spaceFileName string, flagValue *string, api client.SpacesClient) (*spaces.Space, error) {
	providedSpaceID := ""
	if flagValue != nil && *flagValue != "" {
		providedSpaceID = *flagValue
	} else if args != nil && len(args) >= 1 {
		providedSpaceID = args[0]
	}

	if providedSpaceID != "" {
		teamID := viper.GetString("teamID")
		if teamID != "" {
			return api.Spaces().GetByTeamAndName(teamID, providedSpaceID)
		} else {
			return api.Spaces().GetByID(providedSpaceID)
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
			errors.New("missing space ID (--space, -s) or"),
			fmt.Errorf("no spacefile found at %s or", spaceFileName),
			errors.New("missing team ID (--team, -t or $SPACES_TEAM_ID)"),
		)
		return nil, err
	}

	return nil, err
}
