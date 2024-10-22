package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetMode(t *testing.T) {
	router := setupRouter()

	mode.Enabled = true

	// fetch the original mode using the API endpoint
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/mode", nil)
	router.ServeHTTP(w, req)

	// assert the API endpoint returns the mode correctly
	assert.Equal(t, http.StatusOK, w.Code)
	require.JSONEq(t, `{"enabled":true}`, w.Body.String())

	// fetch the update mode using the API endpoint
	mode.Enabled = false
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/api/mode", nil)
	router.ServeHTTP(w, req)

	// assert the API endpoint returns the adjusted mode correctly
	assert.Equal(t, http.StatusOK, w.Code)
	require.JSONEq(t, `{"enabled":false}`, w.Body.String())
}

func TestUpdateMode(t *testing.T) {
	router := setupRouter()

	mode.Enabled = true

	// update the mode using the API endpoint
	w := httptest.NewRecorder()
	body := `{"enabled":false}`
	req, _ := http.NewRequest(http.MethodPut, "/api/mode", strings.NewReader(body))
	router.ServeHTTP(w, req)

	// make sure the API endpoint responds correctly
	assert.Equal(t, http.StatusOK, w.Code)
	require.JSONEq(t, `{"enabled":false}`, w.Body.String())

	// make sure the mode was updated properly
	assert.Equal(t, false, mode.Enabled)
}

func TestUpdateModeWithInvalidBody(t *testing.T) {
	router := setupRouter()

	mode.Enabled = true

	// try to update the mode using the API endpoint with an invalid body
	w := httptest.NewRecorder()
	body := `{"foo":"bar}`
	req, _ := http.NewRequest(http.MethodPut, "/api/mode", strings.NewReader(body))
	router.ServeHTTP(w, req)

	// make sure the API endpoint responds correctly
	assert.Equal(t, http.StatusBadRequest, w.Code)
	require.JSONEq(t, `{"error": "Invalid request body"}`, w.Body.String())

	// make sure the mode not updated
	assert.Equal(t, true, mode.Enabled)
}
