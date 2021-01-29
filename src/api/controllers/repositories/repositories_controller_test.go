package repositories

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/PS-07/go-microservices/src/api/clients/rest"
	"github.com/PS-07/go-microservices/src/api/domain/repositories"
	"github.com/PS-07/go-microservices/src/api/utils/errors"
	"github.com/PS-07/go-microservices/src/api/utils/testutils"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	rest.StartMockup()
	os.Exit(m.Run())
}

func TestCreateRepoInvalidJsonRequest(t *testing.T) {
	request, _ := http.NewRequest(http.MethodPost, "/repositories", strings.NewReader(``))
	response := httptest.NewRecorder()
	c := testutils.GetMockedContext(request, response)

	CreateRepo(c)
	assert.EqualValues(t, http.StatusBadRequest, response.Code)

	apiErr, err := errors.NewAPIErrFromBytes(response.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusBadRequest, apiErr.StatusFn())
	assert.EqualValues(t, "invalid json body", apiErr.MessageFn())
}

func TestCreateRepoErrorFromGithub(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockup(rest.Mock{
		URL:        "https://api.github.com/user/repos",
		HTTPMethod: http.MethodPost,
		Response: &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body: ioutil.NopCloser(strings.NewReader(`{
				"message": "Requires authentication",
				"documentation_url": "https://docs.github.com/rest/reference/repos#create-a-repository-for-the-authenticated-user"
			}`)),
		},
	})

	request, _ := http.NewRequest(http.MethodPost, "/repositories", strings.NewReader(`{"name": "test-repo"}`))
	response := httptest.NewRecorder()
	c := testutils.GetMockedContext(request, response)

	CreateRepo(c)
	assert.EqualValues(t, http.StatusUnauthorized, response.Code)

	apiErr, err := errors.NewAPIErrFromBytes(response.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusUnauthorized, apiErr.StatusFn())
	assert.EqualValues(t, "Requires authentication", apiErr.MessageFn())
}

func TestCreateRepoNoError(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockup(rest.Mock{
		URL:        "https://api.github.com/user/repos",
		HTTPMethod: http.MethodPost,
		Response: &http.Response{
			StatusCode: http.StatusCreated,
			Body:       ioutil.NopCloser(strings.NewReader(`{"id": 123}`)),
		},
	})

	request, _ := http.NewRequest(http.MethodPost, "/repositories", strings.NewReader(`{"name": "test-repo"}`))
	response := httptest.NewRecorder()
	c := testutils.GetMockedContext(request, response)

	CreateRepo(c)
	assert.EqualValues(t, http.StatusCreated, response.Code)

	var result repositories.CreateRepoResponse
	err := json.Unmarshal(response.Body.Bytes(), &result)
	assert.Nil(t, err)
	assert.EqualValues(t, 123, result.ID)
	assert.EqualValues(t, "", result.Name)
	assert.EqualValues(t, "", result.Owner)
}
