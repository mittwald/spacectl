package spacefile

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewSoftwareIsAdded(t *testing.T) {
	base := SoftwareDefList{SoftwareDef{Identifier: "typo3", Version: "1.2.3"}}
	overlay := SoftwareDefList{SoftwareDef{Identifier: "magento", Version: "3.2.1"}}

	merged, err := base.Merge(overlay)

	require.Nil(t, err)
	require.Len(t, merged, 2)
	require.Equal(t, "typo3", merged[0].Identifier)
	require.Equal(t, "magento", merged[1].Identifier)
}

func TestExistingSoftwareIsMerged(t *testing.T) {
	base := SoftwareDefList{SoftwareDef{Identifier: "typo3", Version: "1.2.3"}}
	overlay := SoftwareDefList{SoftwareDef{Identifier: "typo3", Version: "3.2.1"}}

	merged, err := base.Merge(overlay)

	require.Nil(t, err)
	require.Len(t, merged, 1)
	require.Equal(t, "typo3", merged[0].Identifier)
	require.Equal(t, "3.2.1", merged[0].Version)
}