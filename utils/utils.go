package utils

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

type Route struct {
	ID        string `csv:"route_id"`
	AgencyID  string `csv:"agency_id"`
	ShortName string `csv:"route_short_name"`
	LongName  string `csv:"route_long_name"`
	Desc      string `csv:"route_desc"`
	Type      string `csv:"route_type"`
}

func GetRoutes() []Route {
	var routes []Route
	path := os.Getenv("ROUTES_FILE")
	if path == "" {
		path = "../../static/routesinfo/routes.txt"
	}
	csvFile, err := os.Open(path)
	if err != nil {
		panic(err)
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
			log.Fatal(err)
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

	return routes
}

func SearchLine(line string, routes []Route) []Route {
	result := []Route{}
	re := regexp.MustCompile("[0-9]+")

	// Try to find the line trimming the letters in the ShortName
	for _, route := range routes {
		lineNumber := re.FindString(route.ShortName)
		if ok := lineNumber == line; ok {
			result = append(result, route)
		}
	}

	// In case of no results, try to find the line contained in the ShortName
	if len(result) == 0 {
		for _, route := range routes {
			if ok := strings.Contains(route.ShortName, line); ok {
				result = append(result, route)
			}
		}
	}
	return result
}
