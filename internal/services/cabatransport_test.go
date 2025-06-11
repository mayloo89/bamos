package services

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// *** NewAPIClient tests ***

func TestNewAPIClient(t *testing.T) {
	// Call NewAPIClient to create a new API client
	client := NewAPIClient()

	// Assertions
	require.NotNil(t, client)
	assert.Equal(t, BaseURL, client.BaseURL)

	// Check the HTTP client
	require.NotNil(t, client.HTTPClient)
	httpClient, ok := client.HTTPClient.(*http.Client)
	require.True(t, ok)
	assert.Equal(t, DefaultTimeout, httpClient.Timeout)
}

// *** ParkingRules tests ***

func TestParkingRules_OK(t *testing.T) {
	// Create a mock HTTP client
	mockClient := new(MockAPIClient)

	// Simulate a successful API response
	mockResponseBody := `{
        "totalFull": 1,
        "instancias": [
            {
                "nombre": "Test Rule",
                "claseId": "1",
                "clase": "Test Class",
                "id": "123",
                "distancia": "100",
                "contenido": {
                    "contenido": [
                        {
                            "nombreId": "calle",
                            "nombre": "Calle",
                            "position": "1",
                            "valor": "Corrientes"
                        },
                        {
                            "nombreId": "altura",
                            "nombre": "Altura",
                            "position": "2",
                            "valor": "1000"
                        },
                        {
                            "nombreId": "permiso",
                            "nombre": "Permiso",
                            "position": "3",
                            "valor": "Permitido"
                        },
                        {
                            "nombreId": "horario",
                            "nombre": "Horario",
                            "position": "4",
                            "valor": "08:00-20:00"
                        },
                        {
                            "nombreId": "lado",
                            "nombre": "Lado",
                            "position": "5",
                            "valor": "Izquierdo"
                        }
                    ]
                }
            }
        ],
        "total": 1
    }`

	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(mockResponseBody)),
		Header:     make(http.Header),
	}

	// Set up the mock to return the simulated response
	mockClient.On("Do", mock.Anything).Return(mockResponse, nil)

	// Create the API client with the mock HTTP client
	apiClient := &Client{
		BaseURL:    BaseURL,
		HTTPClient: mockClient,
	}

	// Call the ParkingRules method
	lat, long := -34.603722, -58.381592
	rules, err := apiClient.ParkingRules(lat, long)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, rules)
	assert.Contains(t, rules, "Corrientes 1000")
	assert.Equal(t, []string{"Lado izquierdo: permitido las 08:00-20:00."}, rules["Corrientes 1000"])

	// Verify that the mock was called as expected
	mockClient.AssertExpectations(t)
}

func TestParkingRules_RequestError(t *testing.T) {
	// Create a mock HTTP client
	mockClient := new(MockAPIClient)

	// Simulate a request creation error by returning nil for the response
	mockClient.On("Do", mock.Anything).Return(nil, errors.New("request error"))

	// Create the API client with the mock HTTP client
	apiClient := &Client{
		BaseURL:    BaseURL,
		HTTPClient: mockClient,
	}

	// Call the ParkingRules method
	lat, long := -34.603722, -58.381592
	rules, err := apiClient.ParkingRules(lat, long)

	// Assertions
	require.Error(t, err)
	assert.Nil(t, rules)
	assert.Equal(t, "no parking rules found", err.Error())

	// Verify that the mock was called as expected
	mockClient.AssertExpectations(t)
}

func TestParkingRules_ResponseError(t *testing.T) {
	// Create a mock HTTP client
	mockClient := new(MockAPIClient)

	// Simulate a non-200 response
	mockResponse := &http.Response{
		StatusCode: http.StatusInternalServerError,
		Body:       io.NopCloser(strings.NewReader("")),
		Header:     make(http.Header),
	}
	mockClient.On("Do", mock.Anything).Return(mockResponse, nil)

	// Create the API client with the mock HTTP client
	apiClient := &Client{
		BaseURL:    BaseURL,
		HTTPClient: mockClient,
	}

	// Call the ParkingRules method
	lat, long := -34.603722, -58.381592
	rules, err := apiClient.ParkingRules(lat, long)

	// Assertions
	require.Error(t, err)
	assert.Nil(t, rules)

	// Verify that the mock was called as expected
	mockClient.AssertExpectations(t)
}

func TestParkingRules_NoRulesFound(t *testing.T) {
	// Create a mock HTTP client
	mockClient := new(MockAPIClient)

	// Simulate an API response with no rules
	mockResponseBody := `{
        "totalFull": 0,
        "instancias": [],
        "total": 0
    }`
	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(mockResponseBody)),
		Header:     make(http.Header),
	}
	mockClient.On("Do", mock.Anything).Return(mockResponse, nil)

	// Create the API client with the mock HTTP client
	apiClient := &Client{
		BaseURL:    BaseURL,
		HTTPClient: mockClient,
	}

	// Call the ParkingRules method
	lat, long := -34.603722, -58.381592
	rules, err := apiClient.ParkingRules(lat, long)

	// Assertions
	require.Error(t, err)
	assert.Nil(t, rules)
	assert.Equal(t, "no parking rules found", err.Error())

	// Verify that the mock was called as expected
	mockClient.AssertExpectations(t)
}

func TestParkingRules_InvalidJSON(t *testing.T) {
	// Create a mock HTTP client
	mockClient := new(MockAPIClient)

	// Simulate an API response with invalid JSON
	mockResponseBody := `invalid json`
	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(mockResponseBody)),
		Header:     make(http.Header),
	}
	mockClient.On("Do", mock.Anything).Return(mockResponse, nil)

	// Create the API client with the mock HTTP client
	apiClient := &Client{
		BaseURL:    BaseURL,
		HTTPClient: mockClient,
	}

	// Call the ParkingRules method
	lat, long := -34.603722, -58.381592
	rules, err := apiClient.ParkingRules(lat, long)

	// Assertions
	require.Error(t, err)
	assert.Nil(t, rules)

	// Verify that the mock was called as expected
	mockClient.AssertExpectations(t)
}

func TestParkingRules_RetryLogic(t *testing.T) {
	// Create a mock HTTP client
	mockClient := new(MockAPIClient)

	// Simulate multiple failed attempts before success
	mockResponseBody := `{
        "totalFull": 1,
        "instancias": [
            {
                "nombre": "Test Rule",
                "claseId": "1",
                "clase": "Test Class",
                "id": "123",
                "distancia": "100",
                "contenido": {
                    "contenido": [
                        {
                            "nombreId": "calle",
                            "nombre": "Calle",
                            "position": "1",
                            "valor": "Corrientes"
                        },
                        {
                            "nombreId": "altura",
                            "nombre": "Altura",
                            "position": "2",
                            "valor": "1000"
                        },
                        {
                            "nombreId": "permiso",
                            "nombre": "Permiso",
                            "position": "3",
                            "valor": "Permitido"
                        },
                        {
                            "nombreId": "horario",
                            "nombre": "Horario",
                            "position": "4",
                            "valor": "08:00-20:00"
                        },
                        {
                            "nombreId": "lado",
                            "nombre": "Lado",
                            "position": "5",
                            "valor": "Izquierdo"
                        }
                    ]
                }
            }
        ],
        "total": 1
    }`
	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(mockResponseBody)),
		Header:     make(http.Header),
	}

	// First two attempts fail, third succeeds
	mockClient.On("Do", mock.Anything).Return(nil, errors.New("temporary error")).Once()
	mockClient.On("Do", mock.Anything).Return(nil, errors.New("temporary error")).Once()
	mockClient.On("Do", mock.Anything).Return(mockResponse, nil).Once()

	// Create the API client with the mock HTTP client
	apiClient := &Client{
		BaseURL:    BaseURL,
		HTTPClient: mockClient,
	}

	// Call the ParkingRules method
	lat, long := -34.603722, -58.381592
	rules, err := apiClient.ParkingRules(lat, long)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, rules)
	assert.Contains(t, rules, "Corrientes 1000")
	assert.Equal(t, []string{"Lado izquierdo: permitido las 08:00-20:00."}, rules["Corrientes 1000"])

	// Verify that the mock was called as expected
	mockClient.AssertExpectations(t)
}

func TestParkingRules_NilResponse(t *testing.T) {
	// Create a mock HTTP client
	mockClient := new(MockAPIClient)

	// Simulate a nil response
	mockClient.On("Do", mock.Anything).Return(nil, nil)

	// Create the API client with the mock HTTP client
	apiClient := &Client{
		BaseURL:    BaseURL,
		HTTPClient: mockClient,
	}

	// Call the ParkingRules method
	lat, long := -34.603722, -58.381592
	rules, err := apiClient.ParkingRules(lat, long)

	// Assertions
	require.Error(t, err)
	assert.Nil(t, rules)
	assert.Equal(t, "no parking rules found", err.Error())

	// Verify that the mock was called as expected
	mockClient.AssertExpectations(t)
}

func TestParkingRules_NilResponseBody(t *testing.T) {
	// Create a mock HTTP client
	mockClient := new(MockAPIClient)

	// Simulate a response with a nil body
	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       nil,
		Header:     make(http.Header),
	}
	mockClient.On("Do", mock.Anything).Return(mockResponse, nil)

	// Create the API client with the mock HTTP client
	apiClient := &Client{
		BaseURL:    BaseURL,
		HTTPClient: mockClient,
	}

	// Call the ParkingRules method
	lat, long := -34.603722, -58.381592
	rules, err := apiClient.ParkingRules(lat, long)

	// Assertions
	require.Error(t, err)
	assert.Nil(t, rules)
	assert.Equal(t, "no parking rules found", err.Error())

	// Verify that the mock was called as expected
	mockClient.AssertExpectations(t)
}

func TestParkingRules_RetryLogicWithNilResponse(t *testing.T) {
	// Create a mock HTTP client
	mockClient := new(MockAPIClient)

	// Simulate multiple failed attempts with nil responses
	mockClient.On("Do", mock.Anything).Return(nil, errors.New("temporary error")).Times(DefaultRetries)

	// Create the API client with the mock HTTP client
	apiClient := &Client{
		BaseURL:    BaseURL,
		HTTPClient: mockClient,
	}

	// Call the ParkingRules method
	lat, long := -34.603722, -58.381592
	rules, err := apiClient.ParkingRules(lat, long)

	// Assertions
	require.Error(t, err)
	assert.Nil(t, rules)
	assert.Contains(t, err.Error(), "no parking rules found")

	// Verify that the mock was called as expected
	mockClient.AssertExpectations(t)
}

func TestParkingRules_NewRequestError(t *testing.T) {
	// Create a mock HTTP client
	mockClient := new(MockAPIClient)

	// Create the API client with the mock HTTP client
	apiClient := &Client{
		BaseURL:    "://invalid-url", // Invalid URL to trigger NewRequest error
		HTTPClient: mockClient,
	}

	// Call the ParkingRules method
	lat, long := -34.603722, -58.381592
	rules, err := apiClient.ParkingRules(lat, long)

	// Assertions
	require.Error(t, err)
	assert.Nil(t, rules)
	assert.Contains(t, err.Error(), "error creating request")
}

func TestParkingRules_ReadBodyError(t *testing.T) {
	// Create a mock HTTP client
	mockClient := new(MockAPIClient)

	// Simulate a response with a faulty body reader
	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(&faultyReader{}), // Custom faulty reader
		Header:     make(http.Header),
	}
	mockClient.On("Do", mock.Anything).Return(mockResponse, nil)

	// Create the API client with the mock HTTP client
	apiClient := &Client{
		BaseURL:    BaseURL,
		HTTPClient: mockClient,
	}

	// Call the ParkingRules method
	lat, long := -34.603722, -58.381592
	rules, err := apiClient.ParkingRules(lat, long)

	// Assertions
	require.Error(t, err)
	assert.Nil(t, rules)
	assert.Contains(t, err.Error(), "error reading response body")

	// Verify that the mock was called as expected
	mockClient.AssertExpectations(t)
}

func TestParkingRules_OKWithParidad(t *testing.T) {
	// Create a mock HTTP client
	mockClient := new(MockAPIClient)

	// Simulate a successful API response with "paridad"
	mockResponseBody := `{
        "totalFull": 1,
        "instancias": [
            {
                "nombre": "Test Rule",
                "claseId": "1",
                "clase": "Test Class",
                "id": "123",
                "distancia": "100",
                "contenido": {
                    "contenido": [
                        {
                            "nombreId": "calle",
                            "nombre": "Calle",
                            "position": "1",
                            "valor": "Corrientes"
                        },
                        {
                            "nombreId": "altura",
                            "nombre": "Altura",
                            "position": "2",
                            "valor": "1000"
                        },
                        {
                            "nombreId": "permiso",
                            "nombre": "Permiso",
                            "position": "3",
                            "valor": "Permitido"
                        },
                        {
                            "nombreId": "horario",
                            "nombre": "Horario",
                            "position": "4",
                            "valor": "08:00-20:00"
                        },
                        {
                            "nombreId": "lado",
                            "nombre": "Lado",
                            "position": "5",
                            "valor": "Izquierdo"
                        },
                        {
                            "nombreId": "paridad",
                            "nombre": "Paridad",
                            "position": "6",
                            "valor": "Impar"
                        },
                        {
                            "nombreId": "unexpected",
                            "nombre": "unexpected",
                            "position": "7",
                            "valor": "something"
                        }
                    ]
                }
            }
        ],
        "total": 1
    }`

	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(mockResponseBody)),
		Header:     make(http.Header),
	}

	// Set up the mock to return the simulated response
	mockClient.On("Do", mock.Anything).Return(mockResponse, nil)

	// Create the API client with the mock HTTP client
	apiClient := &Client{
		BaseURL:    BaseURL,
		HTTPClient: mockClient,
	}

	// Call the ParkingRules method
	lat, long := -34.603722, -58.381592
	rules, err := apiClient.ParkingRules(lat, long)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, rules)
	assert.Contains(t, rules, "Corrientes 1000")
	assert.Equal(t, []string{"Lado izquierdo (impar): permitido las 08:00-20:00."}, rules["Corrientes 1000"])

	// Verify that the mock was called as expected
	mockClient.AssertExpectations(t)
}

// *** Helper functions ***

// Custom faulty reader to simulate io.ReadAll error
type faultyReader struct{}

func (f *faultyReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("error reading response body")
}
func (f *faultyReader) Close() error {
	return nil
}
