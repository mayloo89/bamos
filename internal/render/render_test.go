package render

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mayloo89/bamos/internal/model"
)

func Test_AddDefaultData_Success(t *testing.T) {
	assert := assert.New(t)
	required := require.New(t)

	var td model.TemplateData

	r, err := getTestSession()
	required.Nil(err)

	result := AddDefaultData(&td, r)

	assert.NotNil(result)
	assert.NotNil(result.CSRFToken)
}

func Test_RenderTemplate_Success(t *testing.T) {
	assert := assert.New(t)
	required := require.New(t)

	tc, err := CreateTemplateCache()
	required.Nil(err)

	app.UseCache = true
	app.TemplateCache = tc

	var td model.TemplateData
	r, err := getTestSession()
	required.Nil(err)

	var ww testWriter

	err = RenderTemplate(&ww, r, "home.page.tmpl", &td)

	assert.Nil(err)
}

func Test_RenderTemplate_Error(t *testing.T) {
	assert := assert.New(t)
	required := require.New(t)

	tc, err := CreateTemplateCache()
	required.Nil(err)

	app.UseCache = true
	app.TemplateCache = tc

	var td model.TemplateData
	r, err := getTestSession()
	required.Nil(err)

	var ww testWriter

	err = RenderTemplate(&ww, r, "non-existent.page.tmpl", &td)

	assert.NotNil(err)
	assert.Equal("Template non-existent.page.tmpl does not exist in cache (len 3)", err.Error())
}

func Test_NewTemplate(t *testing.T) {
	NewTemplates(app)
}

func Test_CreateTemplateCache(t *testing.T) {
	assert := assert.New(t)
	required := require.New(t)

	tc, err := CreateTemplateCache()
	required.Nil(err)

	assert.Equal(3, len(tc))
}

func getTestSession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/test-url", nil)
	if err != nil {
		return nil, err
	}

	ctx := r.Context()
	ctx, err = session.Load(ctx, r.Header.Get("X-Session"))
	if err != nil {
		return nil, err
	}

	r = r.WithContext(ctx)
	return r, nil
}
