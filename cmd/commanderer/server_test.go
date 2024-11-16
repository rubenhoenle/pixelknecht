package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetServer(t *testing.T) {
	router := setupRouter()

	server.Host = "1.1.1.1"
	server.Port = 1234

	// fetch the original pixelflut server configuration using the API endpoint
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/server", nil)
	router.ServeHTTP(w, req)

	// assert the API endpoint returns the pixelflut server configuration correctly
	assert.Equal(t, http.StatusOK, w.Code)
	require.JSONEq(t, `{"host":"1.1.1.1","port":1234}`, w.Body.String())

	// update the pixelflut server configuration directly
	server.Host = "9.9.9.9"
	server.Port = 4321

	// fetch the updated pixelflut server configuration using the API endpoint
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/api/server", nil)
	router.ServeHTTP(w, req)

	// assert the API endpoint returns the adjusted pixelflut server configuration correctly
	assert.Equal(t, http.StatusOK, w.Code)
	require.JSONEq(t, `{"host":"9.9.9.9","port":4321}`, w.Body.String())
}

func TestUpdateServer(t *testing.T) {
	router := setupRouter()

	server.Host = "1.1.1.1"
	server.Port = 1234

	// update the pixelflut server configuration using the API endpoint
	w := httptest.NewRecorder()
	body := `{"host":"9.9.9.9","port":4321}`
	req, _ := http.NewRequest(http.MethodPut, "/api/server", strings.NewReader(body))
	router.ServeHTTP(w, req)

	// make sure the API endpoint responds correctly
	assert.Equal(t, http.StatusOK, w.Code)
	require.JSONEq(t, `{"host":"9.9.9.9","port":4321}`, w.Body.String())

	// make sure the pixelflut server configuration was updated properly
	assert.Equal(t, "9.9.9.9", server.Host)
	assert.Equal(t, 4321, server.Port)
}

func TestUpdateServerWithInvalidBody(t *testing.T) {
	router := setupRouter()

	server.Host = "1.1.1.1"
	server.Port = 1234

	// try to update the pixelflut server configuration using the API endpoint with an invalid body
	w := httptest.NewRecorder()
	body := `{"foo":"bar}`
	req, _ := http.NewRequest(http.MethodPut, "/api/server", strings.NewReader(body))
	router.ServeHTTP(w, req)

	// make sure the API endpoint responds correctly
	assert.Equal(t, http.StatusBadRequest, w.Code)
	require.JSONEq(t, `{"error": "Invalid request body"}`, w.Body.String())

	// make sure the pixelflut server configuration was not updated
	assert.Equal(t, "1.1.1.1", server.Host)
	assert.Equal(t, 1234, server.Port)
}
