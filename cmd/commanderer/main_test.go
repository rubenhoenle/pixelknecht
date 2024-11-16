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
	mode.PosY = 45
	mode.PosX = 10
	mode.ScaleFactor = float64(1)
	mode.ImageUrl = "https://test.com/image.jpg"

	// fetch the original mode using the API endpoint
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/mode", nil)
	router.ServeHTTP(w, req)

	// assert the API endpoint returns the mode correctly
	assert.Equal(t, http.StatusOK, w.Code)
	require.JSONEq(t, `{"enabled":true,"posY":45,"posX":10,"scaleFactor":1,"imageUrl":"https://test.com/image.jpg"}`, w.Body.String())

	// update the mode directly
	mode.Enabled = false
	mode.PosY = 55
	mode.PosX = 5
	mode.ScaleFactor = float64(2)
	mode.ImageUrl = "https://example.com/image.png"

	// fetch the updated mode using the API endpoint
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/api/mode", nil)
	router.ServeHTTP(w, req)

	// assert the API endpoint returns the adjusted mode correctly
	assert.Equal(t, http.StatusOK, w.Code)
	require.JSONEq(t, `{"enabled":false,"posY":55,"posX":5,"scaleFactor":2,"imageUrl":"https://example.com/image.png"}`, w.Body.String())
}

func TestUpdateMode(t *testing.T) {
	router := setupRouter()

	mode.Enabled = true
	mode.PosY = 45
	mode.PosX = 10
	mode.ScaleFactor = float64(1)
	mode.ImageUrl = "https://test.com/image.jpg"

	// update the mode using the API endpoint
	w := httptest.NewRecorder()
	body := `{"enabled":false,"posY":15,"posX":20,"scaleFactor":2,"imageUrl":"https://example.com/image.png"}`
	req, _ := http.NewRequest(http.MethodPut, "/api/mode", strings.NewReader(body))
	router.ServeHTTP(w, req)

	// make sure the API endpoint responds correctly
	assert.Equal(t, http.StatusOK, w.Code)
	require.JSONEq(t, `{"enabled":false,"posY":15,"posX":20,"scaleFactor":2,"imageUrl":"https://example.com/image.png"}`, w.Body.String())

	// make sure the mode was updated properly
	assert.Equal(t, false, mode.Enabled)
	assert.Equal(t, 15, mode.PosY)
	assert.Equal(t, 20, mode.PosX)
	assert.Equal(t, float64(2), mode.ScaleFactor)
	assert.Equal(t, "https://example.com/image.png", mode.ImageUrl)
}

func TestUpdateModeWithInvalidBody(t *testing.T) {
	router := setupRouter()

	mode.Enabled = true
	mode.PosY = 45
	mode.PosX = 10
	mode.ScaleFactor = float64(1)
	mode.ImageUrl = "https://test.com/image.jpg"

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
	assert.Equal(t, 45, mode.PosY)
	assert.Equal(t, 10, mode.PosX)
	assert.Equal(t, float64(1), mode.ScaleFactor)
	assert.Equal(t, "https://test.com/image.jpg", mode.ImageUrl)
}
