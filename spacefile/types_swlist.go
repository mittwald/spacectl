package spacefile

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/viper"

	"github.com/Masterminds/semver"
	"github.com/imdario/mergo"
	"github.com/mittwald/spacectl/client"
)

type SoftwareDef struct {
	Identifier string      `hcl:",key"`
	Version    string      `hcl:"version"`
	UserData   interface{} `hcl:"userData"`
}

type SoftwareDefList []SoftwareDef

func (l SoftwareDefList) copy() SoftwareDefList {
	c := make(SoftwareDefList, len(l))
	for i := range l {
		c[i] = l[i]
	}
	return c
}

func (l SoftwareDefList) Merge(other SoftwareDefList) (SoftwareDefList, error) {
	merged := l.copy()

	for i := range other {
		foundMatch := -1
		for j := range merged {
			if other[i].Identifier == merged[j].Identifier {
				foundMatch = j
				break
			}
		}

		if foundMatch >= 0 {
			err := mergo.Merge(&merged[foundMatch], &other[i], mergo.WithOverride)
			if err != nil {
				return nil, err
			}
		} else {
			merged = append(merged, other[i])
		}
	}

	return merged, nil
}

// Validate performs (optional) online validation of software version and name
func (s SoftwareDef) Validate(offline bool) error {
	constrait, errSem := semver.NewConstraint(s.Version)
	if errSem != nil {
		return fmt.Errorf("version: %s", errSem.Error())
	}

	if !offline {
		api, err := client.NewSpacesClient(client.SpacesClientConfig{
			APIServer: viper.GetString("apiServer"),
			Token:     viper.GetString("token"),
		})
		if err != nil {
			return err
		}

		appList, err := api.Applications().List()
		if err != nil {
			return err
		}

		for _, app := range appList {
			if app.Identifier == s.Identifier {
				sw, err := api.Applications().Get(app.Identifier)
				if err != nil {
					return err
				}

				for _, v := range sw.Versions {
					if semVer, err := semver.NewVersion(v.Number); err == nil {
						if constrait.Check(semVer) {
							return nil
						}
					}
				}

				// fancy error
				versionHelp := color.BlueString("use ") + color.YellowString("spacectl software apps show %s", s.Identifier) + color.BlueString(" to list available applications")
				return fmt.Errorf("version %s is not available\n\n%s", s.Version, versionHelp)
			}
		}
		// fancy error 2
		softwareHelp := color.BlueString("use ") + color.YellowString("spacectl software apps list") + color.BlueString(" to list available applications")
		return fmt.Errorf("software %s is not available\n\n%s", s.Identifier, softwareHelp)
	}
	return nil
}
