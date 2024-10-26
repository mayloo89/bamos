package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_run_Success(t *testing.T) {
	assert := assert.New(t)

	err := run()

	assert.Nil(err)
}
