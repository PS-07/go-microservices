package github

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateRepoRequestAsJson(t *testing.T) {
	request := CreateRepoRequest{
		Name:        "golang introduction",
		Description: "a golang introduction repository",
		Homepage:    "https://github.com",
		Private:     false,
		HasIssues:   true,
		HasProjects: true,
		HasWiki:     true,
	}

	jsonBytes, err := json.Marshal(request)

	assert.Nil(t, err)
	assert.NotNil(t, jsonBytes)

	var target CreateRepoRequest
	err = json.Unmarshal(jsonBytes, &target)

	assert.Nil(t, err)
	assert.EqualValues(t, target.Name, request.Name)
	assert.EqualValues(t, target.HasIssues, request.HasIssues)
}
