package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_run_Success(t *testing.T) {
	assert := assert.New(t)

	db, err := run()

	assert.Nil(err)
	assert.NotNil(db)
}
