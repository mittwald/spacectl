package payment

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSpaceResourcePreprovisioningInputJSONSerializationWithFullResources(t *testing.T) {
	r := SpaceResourcePreprovisioningInput{
		Storage: &SpaceResourcePreprovisioningInputItem{Quantity: 1024},
		Stages: &SpaceResourcePreprovisioningInputItem{Quantity: 3},
		Scaling: &SpaceResourcePreprovisioningInputItem{Quantity: 8},
	}

	b, err := json.Marshal(&r)
	if err != nil {
		t.Fatal(err)
	}

	target := make(map[string]interface{})
	err = json.Unmarshal(b, &target)
	if err != nil {
		t.Fatal(err)
	}

	_, hasStorage := target["storage"]
	_, hasStages := target["stages"]
	_, hasScaling := target["scaling"]

	assert.True(t, hasStorage)
	assert.True(t, hasStages)
	assert.True(t, hasScaling)
}

func TestSpaceResourcePreprovisioningInputJSONSerializationWithSparseResources(t *testing.T) {
	r := SpaceResourcePreprovisioningInput{
		Storage: &SpaceResourcePreprovisioningInputItem{
			Quantity: 1024,
		},
	}

	b, err := json.Marshal(&r)
	if err != nil {
		t.Fatal(err)
	}

	target := make(map[string]interface{})
	err = json.Unmarshal(b, &target)
	if err != nil {
		t.Fatal(err)
	}

	_, hasStorage := target["storage"]
	_, hasStages := target["stages"]
	_, hasScaling := target["scaling"]

	assert.True(t, hasStorage)
	assert.False(t, hasStages)
	assert.False(t, hasScaling)
}