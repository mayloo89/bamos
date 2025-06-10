package main

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayloo89/bamos/internal/config"
	"github.com/mayloo89/bamos/internal/handler"
)

func Test_routes_Success(t *testing.T) {
	assert := assert.New(t)
	ac := &config.AppConfig{}
	repo := handler.NewRepo(ac, nil)

	handler := routes(ac, repo)

	var expectedType bool
	switch handler.(type) {
	case http.Handler:
		expectedType = true
	default:
		expectedType = false
	}

	assert.NotNil(handler)
	assert.True(expectedType)
}
