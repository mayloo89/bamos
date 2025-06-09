package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_run_Success(t *testing.T) {
	assert := assert.New(t)

	// Set the ROUTES_FILE env var so utils.GetRoutes works regardless of test working dir
	t.Setenv("ROUTES_FILE", "../../static/routesinfo/routes.txt")

	err := run()

	assert.Nil(err, "run() should not return an error when ROUTES_FILE is set and file exists")
}
