package spacefile

import "testing"
import "github.com/stretchr/testify/require"

func errMsg(e error) string {
	if e != nil {
		return e.Error()
	}
	return ""
}

func TestSpacefileCanBeCorrectlyParsed(t *testing.T) {
	spacefile, err := ParseSpacefile("./example.hcl")

	require.Nil(t, err, errMsg(err))
	require.NotNil(t, spacefile)
}
