package testutils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMockedContext(t *testing.T) {
	request, err := http.NewRequest(http.MethodGet, "http://localhost:1234/testroute", nil)
	request.Header = http.Header{"X-Mock": {"true"}}
	assert.Nil(t, err)

	response := httptest.NewRecorder()
	c := GetMockedContext(request, response)

	assert.EqualValues(t, http.MethodGet, c.Request.Method)
	assert.EqualValues(t, "1234", c.Request.URL.Port())
	assert.EqualValues(t, "/testroute", c.Request.URL.Path)
	assert.EqualValues(t, "http", c.Request.URL.Scheme)
	assert.EqualValues(t, "true", c.GetHeader("X-Mock"))
	assert.EqualValues(t, 1, len(c.Request.Header))
}
