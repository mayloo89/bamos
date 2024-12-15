package handler

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type postData struct {
	key   string
	value string
}

var tests = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"search GET", "/colectivos/search", "GET", []postData{}, http.StatusOK},
	// TODO: Mock the transport API request for /colectivos/vehiclePositionsSimple
	//{"vehiclePositionsSimple", "/colectivos/vehiclePositionsSimple", "GET", []postData{}, http.StatusOK},
	{"search POST", "/colectivos/search", "POST", []postData{
		{key: "line", value: "25"},
	}, http.StatusOK},
}

func Test_Handlers(t *testing.T) {
	assert := assert.New(t)
	required := require.New(t)
	var resp *http.Response
	var err error

	routes := getRoutes()
	testServer := httptest.NewTLSServer(routes)
	defer testServer.Close()

	for _, test := range tests {
		switch test.method {
		case "GET":
			resp, err = testServer.Client().Get(testServer.URL + test.url)
		case "POST":
			values := url.Values{}
			for _, param := range test.params {
				values.Add(param.key, param.value)
			}
			resp, err = testServer.Client().PostForm(testServer.URL+test.url, values)
		}
		required.Nil(err)
		assert.Equal(test.expectedStatusCode, resp.StatusCode)
	}
}
