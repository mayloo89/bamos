package main

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NoSurf_Success(t *testing.T) {
	assert := assert.New(t)
	var th testHandler
	var expectedType bool

	handler := NoSurf(&th)

	switch handler.(type) {
	case http.Handler:
		expectedType = true
	default:
		expectedType = false
	}

	assert.NotNil(handler)
	assert.True(expectedType)
}

func Test_SessionLoad_Success(t *testing.T) {
	assert := assert.New(t)
	var th testHandler
	var expectedType bool

	handler := SessionLoad(&th)

	switch handler.(type) {
	case http.Handler:
		expectedType = true
	default:
		expectedType = false
	}

	assert.NotNil(handler)
	assert.True(expectedType)
}
