package main

import (
	"net/http"
	"os"
	"testing"
)

type testHandler struct{}

func (th *testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
