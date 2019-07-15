package spaces

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestWithStorageStrSetsStorageToAppropriateByteValue(t *testing.T) {
	input := SpacePaymentLinkInput{}
	err := WithStorageStr("10G")(&input)

	require.Nil(t, err)
	require.NotNil(t, input.Preprovisionings)
	require.Equal(t, uint64(10 * 1024 * 1024 * 1024), input.Preprovisionings.Storage.Quantity)
}

func TestWithStorageStrAllowsZeroValue(t *testing.T) {
	input := SpacePaymentLinkInput{}
	err := WithStorageStr("0G")(&input)

	require.Nil(t, err)

	// this is fine
	if input.Preprovisionings == nil {
		return
	}

	// this is fine, too
	if input.Preprovisionings.Storage == nil {
		return
	}

	require.Equal(t, uint64(0), input.Preprovisionings.Storage.Quantity)
}

func TestWithStorageStrDoesNotAllowNegativeValue(t *testing.T) {
	input := SpacePaymentLinkInput{}
	err := WithStorageStr("-10G")(&input)

	require.NotNil(t, err)
}