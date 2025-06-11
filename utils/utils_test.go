package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func createTempFileWithContent(t *testing.T, content string) string {
	f, err := os.CreateTemp("", "routes_test_*.csv")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if content != "" {
		_, werr := f.WriteString(content)
		if werr != nil {
			t.Fatalf("failed to write to temp file: %v", werr)
		}
	}
	if err := f.Close(); err != nil {
		t.Fatalf("failed to close temp file: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Remove(f.Name()); err != nil {
			t.Errorf("failed to remove temp file: %v", err)
		}
	})
	return f.Name()
}

func TestGetRoutes(t *testing.T) {
	t.Run("valid file", func(t *testing.T) {
		err := os.Setenv("ROUTES_FILE", "../static/routesinfo/routes.txt")
		require.Nil(t, err)

		routes, err := GetRoutes()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(routes) == 0 {
			t.Error("expected some routes, got 0")
		}
	})

	t.Run("file does not exist", func(t *testing.T) {
		err := os.Setenv("ROUTES_FILE", "../static/routesinfo/doesnotexist.txt")
		require.Nil(t, err)

		_, err = GetRoutes()
		if err == nil {
			t.Error("expected error for missing file, got nil")
		}
	})

	t.Run("malformed file", func(t *testing.T) {
		file := createTempFileWithContent(t, "bad,data,onlythree\n")
		err := os.Setenv("ROUTES_FILE", file)
		require.Nil(t, err)

		_, err = GetRoutes()
		if err == nil {
			t.Error("expected error for malformed file, got nil")
		}
	})

	t.Run("extra columns", func(t *testing.T) {
		file := createTempFileWithContent(t, "id,agency,short,long,desc,type,extra1,extra2\n1,2,3,4,5,6,7,8\n")
		err := os.Setenv("ROUTES_FILE", file)
		require.Nil(t, err)

		routes, err := GetRoutes()
		if err != nil {
			t.Fatalf("expected no error for extra columns, got %v", err)
		}
		if len(routes) != 2 {
			t.Errorf("expected 2 routes, got %d", len(routes))
		}
	})

	t.Run("empty file", func(t *testing.T) {
		file := createTempFileWithContent(t, "")
		err := os.Setenv("ROUTES_FILE", file)
		require.Nil(t, err)
		routes, err := GetRoutes()
		if err != nil {
			t.Fatalf("expected no error for empty file, got %v", err)
		}
		if len(routes) != 0 {
			t.Errorf("expected 0 routes for empty file, got %d", len(routes))
		}
	})
}

func TestSearchLine(t *testing.T) {
	routes := []Route{
		{ID: "1", ShortName: "123A"},
		{ID: "2", ShortName: "456B"},
		{ID: "3", ShortName: "123B"},
	}
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"exact numeric match", "123", 2},
		{"partial match", "A", 1},
		{"no match", "999", 0},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := SearchLine(tc.input, routes)
			if len(result) != tc.expected {
				t.Errorf("expected %d results, got %d", tc.expected, len(result))
			}
		})
	}
}
