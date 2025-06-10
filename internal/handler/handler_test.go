package handler

import (
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/mayloo89/bamos/internal/config"
	"github.com/mayloo89/bamos/internal/helpers"
	"github.com/mayloo89/bamos/internal/render"
	"github.com/mayloo89/bamos/internal/services"
)

func setupTestApp(mockAPIClient services.APIClient) (*Repository, *config.AppConfig) {
	app := &config.AppConfig{
		InProduction: false,
		InfoLog:      log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
	}
	repo := NewRepo(app, mockAPIClient)
	render.NewTemplates(app)
	helpers.NewHelpers(app)
	return repo, app
}

func Test_Home(t *testing.T) {
	// Create a mock API client
	mockAPIClient := new(services.MockAPIClient)
	mockAPIClient.On("AllowedParking", mock.Anything, mock.Anything).Return(nil, nil)

	repo, _ := setupTestApp(mockAPIClient)

	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(repo.Home)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func Test_VehiclePositionsSimple(t *testing.T) {
	// Create a mock API client
	mockAPIClient := new(services.MockAPIClient)
	mockAPIClient.On("AllowedParking", mock.Anything, mock.Anything).Return(nil, nil)

	repo, _ := setupTestApp(mockAPIClient)

	req, err := http.NewRequest("GET", "/colectivos/vehiclePositionsSimple", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(repo.VehiclePositionsSimple)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func Test_SearchLine(t *testing.T) {
	// Create a mock API client
	mockAPIClient := new(services.MockAPIClient)
	mockAPIClient.On("AllowedParking", mock.Anything, mock.Anything).Return(nil, nil)

	repo, _ := setupTestApp(mockAPIClient)

	req, err := http.NewRequest("GET", "/colectivos/search", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(repo.SearchLine)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func Test_PostSearchLine_Success(t *testing.T) {
	// Create a mock API client
	mockAPIClient := new(services.MockAPIClient)
	mockAPIClient.On("AllowedParking", mock.Anything, mock.Anything).Return(nil, nil)

	repo, _ := setupTestApp(mockAPIClient)

	form := url.Values{}
	form.Add("line", "test-line")

	req, err := http.NewRequest("POST", "/colectivos/search", strings.NewReader(form.Encode()))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(repo.PostSearchLine)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func Test_AllowedParking(t *testing.T) {
	// Create a mock API client
	mockAPIClient := new(services.MockAPIClient)
	mockAPIClient.On("AllowedParking", mock.Anything, mock.Anything).Return(nil, nil)

	repo, _ := setupTestApp(mockAPIClient)

	req, err := http.NewRequest("GET", "/transit/allowed-parking", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(repo.AllowedParking)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func Test_PostAllowedParking_Success(t *testing.T) {
	// Create a mock API client
	mockAPIClient := new(services.MockAPIClient)
	mockAPIClient.On("ParkingRules", 40.7128, -74.0060).Return(
		services.SimplifiedRules{"Test Rule": {"Detail 1"}}, nil,
	)
	repo, _ := setupTestApp(mockAPIClient)

	form := url.Values{}
	form.Add("latitude", "40.7128")
	form.Add("longitude", "-74.0060")
	form.Add("address", "Test Address")

	req, err := http.NewRequest("POST", "/transit/allowed-parking", strings.NewReader(form.Encode()))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(repo.PostAllowedParking)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func Test_PostAllowedParking_EmptyRules(t *testing.T) {
	// Create a mock API client
	mockAPIClient := new(services.MockAPIClient)
	mockAPIClient.On("ParkingRules", 40.7128, -74.0060).Return(
		services.SimplifiedRules{}, nil,
	)
	repo, _ := setupTestApp(mockAPIClient)

	form := url.Values{}
	form.Add("latitude", "40.7128")
	form.Add("longitude", "-74.0060")
	form.Add("address", "Test Address")

	req, err := http.NewRequest("POST", "/transit/allowed-parking", strings.NewReader(form.Encode()))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(repo.PostAllowedParking)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func Test_PostAllowedParking_ValidationLatitudeError(t *testing.T) {
	// Create a mock API client
	mockAPIClient := new(services.MockAPIClient)
	repo, _ := setupTestApp(mockAPIClient)

	form := url.Values{}
	form.Add("latitude", "ERROR")
	form.Add("longitude", "-74.0060")
	form.Add("address", "Test Address")

	req, err := http.NewRequest("POST", "/transit/allowed-parking", strings.NewReader(form.Encode()))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(repo.PostAllowedParking)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func Test_PostAllowedParking_ValidationLongitudeError(t *testing.T) {
	// Create a mock API client
	mockAPIClient := new(services.MockAPIClient)
	repo, _ := setupTestApp(mockAPIClient)

	form := url.Values{}
	form.Add("latitude", "40.7128")
	form.Add("longitude", "ERROR")
	form.Add("address", "Test Address")

	req, err := http.NewRequest("POST", "/transit/allowed-parking", strings.NewReader(form.Encode()))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(repo.PostAllowedParking)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func Test_PostAllowedParking_ValidationLongitudeEmpty(t *testing.T) {
	// Create a mock API client
	mockAPIClient := new(services.MockAPIClient)
	repo, _ := setupTestApp(mockAPIClient)

	form := url.Values{}
	form.Add("latitude", "40.7128")
	form.Add("address", "Test Address")

	req, err := http.NewRequest("POST", "/transit/allowed-parking", strings.NewReader(form.Encode()))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(repo.PostAllowedParking)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

// func Test_FeedGtfsFrequency(t *testing.T) {
// 	setupTestApp()

// 	req, err := http.NewRequest("GET", "/colectivos/feed-gtfs-frequency", nil)
// 	require.NoError(t, err)

// 	rr := httptest.NewRecorder()
// 	handler := http.HandlerFunc(Repo.FeedGtfsFrequency)

// 	handler.ServeHTTP(rr, req)

// 	assert.Equal(t, http.StatusOK, rr.Code)
// }
