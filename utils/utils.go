package utils

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Route represents a single route entry from the routes CSV file.
type Route struct {
	ID        string `csv:"route_id"`         // Unique route identifier
	AgencyID  string `csv:"agency_id"`        // Agency identifier
	ShortName string `csv:"route_short_name"` // Short name of the route
	LongName  string `csv:"route_long_name"`  // Long name of the route
	Desc      string `csv:"route_desc"`       // Description of the route
	Type      string `csv:"route_type"`       // Type of the route
}

// GetRoutes loads routes from the CSV file specified by the ROUTES_FILE environment variable.
// If the variable is not set, it defaults to static/routesinfo/routes.txt.
// Returns a slice of Route and an error if the file is missing or malformed.
func GetRoutes() ([]Route, error) {
	var routes []Route
	path := os.Getenv("ROUTES_FILE")
	if path == "" {
		path = filepath.Join("static", "routesinfo", "routes.txt")
	}
	csvFile, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open routes file at %s: %w", path, err)
	}
	defer func() {
		if cerr := csvFile.Close(); cerr != nil {
			log.Println("error closing csv file:", cerr)
		}
	}()

	reaader := csv.NewReader(csvFile)
	for {
		line, err := reaader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("error reading routes file at %s: %w", path, err)
		}
		if len(line) < 6 {
			return nil, fmt.Errorf("malformed routes file at %s: expected at least 6 columns, got %d", path, len(line))
		}
		routes = append(routes, Route{
			ID:        line[0],
			AgencyID:  line[1],
			ShortName: line[2],
			LongName:  line[3],
			Desc:      line[4],
			Type:      line[5],
		})
	}

	return routes, nil
}

// SearchLine searches for routes matching the given line string in the provided routes slice.
// It first tries to match by numeric value, then by substring.
func SearchLine(line string, routes []Route) []Route {
	var result []Route
	re := regexp.MustCompile("[0-9]+")

	// First, try to match by numeric value in ShortName
	for _, route := range routes {
		if re.FindString(route.ShortName) == line {
			result = append(result, route)
		}
	}

	// If no results, try substring match in ShortName
	if len(result) == 0 {
		for _, route := range routes {
			if strings.Contains(route.ShortName, line) {
				result = append(result, route)
			}
		}
	}
	return result
}
