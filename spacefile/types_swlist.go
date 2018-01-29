package spacefile

import "github.com/imdario/mergo"

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