package client

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/jarcoal/httpmock.v1"
)

func TestExistingStageProtection(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	clientConfig := SpacesClientConfig{
		Token:     "testingToken",
		APIServer: "https://api.test.spaces.de",
		Logger:    nil,
	}

	var api SpacesClient
	if s, err := NewSpacesClient(clientConfig); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	} else {
		api = s
	}

	spaceID := "22af6ace-2977-4581-8172-cb4c7681954a"
	stageName := "production"

	httpmock.RegisterResponder(
		"GET",
		"https://api.test.spaces.de/v1/spaces/"+spaceID+"/stages/"+stageName+"/protection",
		httpmock.NewStringResponder(200, `{
"type": "oauth",
"modifiedAt": "2018-08-13T13:03:22.965Z"
}`))

	protection, err := api.Spaces().GetStageProtection(spaceID, stageName)
	if err != nil {
		t.Error(err.Error())
	}

	t.Logf("%v", protection)

	assert.Equal(t, "oauth", protection.ProtectionType)
}

func TestNonExistentStageProtection(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	clientConfig := SpacesClientConfig{
		Token:     "testingToken",
		APIServer: "https://api.test.spaces.de",
		Logger:    nil,
	}

	var api SpacesClient
	if s, err := NewSpacesClient(clientConfig); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	} else {
		api = s
	}

	spaceID := "22af6ace-2977-4581-8172-cb4c7681954a"
	stageName := "production"

	httpmock.RegisterResponder(
		"GET",
		"https://api.test.spaces.de/v1/spaces/"+spaceID+"/stages/"+stageName+"/protection",
		httpmock.NewStringResponder(404, ``))

	protection, err := api.Spaces().GetStageProtection(spaceID, stageName)
	if err != nil {
		t.Error(err.Error())
	}

	t.Logf("%v", protection)

	assert.Equal(t, "", protection.ProtectionType)
}
