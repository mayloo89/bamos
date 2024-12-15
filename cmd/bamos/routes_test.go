package main

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayloo89/bamos/internal/config"
)

func Test_routes_Success(t *testing.T) {
	assert := assert.New(t)
	var ac *config.AppConfig
	var expectedType bool

	handler := routes(ac)

	switch handler.(type) {
	case http.Handler:
		expectedType = true
	default:
		expectedType = false
	}

	assert.NotNil(handler)
	assert.True(expectedType)
}
