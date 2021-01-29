package health

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PS-07/go-microservices/src/api/utils/testutils"
	"github.com/stretchr/testify/assert"
)

func TestConstants(t *testing.T) {
	assert.EqualValues(t, "ok", healthStatus)
}

func TestHealth(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	response := httptest.NewRecorder()
	c := testutils.GetMockedContext(request, response)

	Health(c)
	assert.EqualValues(t, http.StatusOK, response.Code)
	assert.EqualValues(t, "ok", response.Body.String())
}
