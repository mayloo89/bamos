package forms

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_New_Success(t *testing.T) {
	assert := assert.New(t)

	form := New(map[string][]string{
		"test-value": {"test"},
	})

	assert.NotNil(form)
	assert.Equal("test", form.Get("test-value"))
}

func Test_Has(t *testing.T) {
	assert := assert.New(t)

	postedValues := url.Values{}
	postedValues.Add("test-value", "test")

	form := New(postedValues)
	ok := form.Has("test-value")
	notOk := form.Has("test-value-2")

	assert.True(ok)
	assert.False(notOk)
}

func Test_Required_Success(t *testing.T) {
	assert := assert.New(t)

	form := New(map[string][]string{
		"test-value": {"test"},
	})

	form.Required("test-value")

	assert.True(form.Valid())
}

func Test_Required_Error(t *testing.T) {
	assert := assert.New(t)

	form := New(map[string][]string{
		"test-value": {"test"},
	})

	form.Required("wrong-field")

	assert.False(form.Valid())
}

func Test_MinLenght_Success(t *testing.T) {
	assert := assert.New(t)

	form := New(map[string][]string{
		"test-value": {"test"},
	})

	form.MinLength("test-value", 2)
	errDetail := form.Errors.Get("test-value")

	assert.True(form.Valid())
	assert.Equal("", errDetail)
}

func Test_MinLenght_Error(t *testing.T) {
	assert := assert.New(t)

	form := New(map[string][]string{
		"test-value": {"test"},
	})

	form.MinLength("test-value", 10)
	errDetail := form.Errors.Get("test-value")

	assert.False(form.Valid())
	assert.Equal("This field must be at least 10 characters long", errDetail)
}
