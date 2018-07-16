package spacefile

import "testing"
import "github.com/stretchr/testify/require"

func errMsg(e error) string {
	if e != nil {
		return e.Error()
	}
	return ""
}

func getParsedSpacefile(t *testing.T) *Spacefile {
	spacefile, err := ParseSpacefile("./example.hcl")

	require.Nil(t, err, errMsg(err))
	require.NotNil(t, spacefile)

	require.Len(t, spacefile.Spaces, 1)
	require.Len(t, spacefile.Spaces[0].Stages, 2)

	return spacefile
}

func TestSpacefileCanBeCorrectlyParsed(t *testing.T) {
	spacefile, err := ParseSpacefile("./example.hcl")

	require.Nil(t, err, errMsg(err))
	require.NotNil(t, spacefile)

	require.Len(t, spacefile.Spaces, 1)
	require.Len(t, spacefile.Spaces[0].Stages, 2)

	stage := spacefile.Spaces[0].Stages[0]

	require.Len(t, stage.Applications, 1)
	require.Len(t, stage.Databases, 1)
}

func TestSpacefileHasCronjobs(t *testing.T) {
	spacefile := getParsedSpacefile(t)
	stage := spacefile.Spaces[0].Stages[0]

	require.Len(t, stage.Cronjobs, 1)

	cron := stage.Cronjobs[0]

	require.Equal(t, "typo3", cron.Identifier)
	require.Equal(t, "*/5 * * * *", cron.Schedule)
	require.Equal(t, true, cron.AllowParallel)
	require.NotNil(t, cron.Command)
	require.Equal(t, "vendor/bin/typo3cmd", cron.Command.Command)
	require.Equal(t, []string{"typo3:scheduler"}, cron.Command.Arguments)
	require.Equal(t, "/var/www/typo3", cron.Command.WorkingDirectory)
}
