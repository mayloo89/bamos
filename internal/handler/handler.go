// Package handler provides HTTP handlers for the Bamos web application.
package handler

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"google.golang.org/protobuf/proto"

	"github.com/mayloo89/bamos/internal/config"
	"github.com/mayloo89/bamos/internal/forms"
	"github.com/mayloo89/bamos/internal/helpers"
	"github.com/mayloo89/bamos/internal/model"
	"github.com/mayloo89/bamos/internal/render"
	"github.com/mayloo89/bamos/internal/services"
	"github.com/mayloo89/bamos/utils"
)

var clientID = "REPLACED"
var clientSecret = "REPLACED"
var apiBaseURL = "https://apitransporte.buenosaires.gob.ar"

type (
	// Repository holds the application config and API client for handlers.
	Repository struct {
		App       *config.AppConfig  // Application configuration
		APIClient services.APIClient // API client for external services
	}
)

// NewRepo creates a new Repository with the given AppConfig and APIClient.
func NewRepo(a *config.AppConfig, apiClient services.APIClient) *Repository {
	return &Repository{
		App:       a,
		APIClient: apiClient,
	}
}

// Home renders the home page.
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {

	stringMap := make(map[string]string)
	stringMap["test"] = "Hello again."

	err := render.RenderTemplate(w, r, "home.page.tmpl", &model.TemplateData{
		StringMap: stringMap,
	})
	if err != nil {
		helpers.ServerError(w, err)
	}
}

// HTTPClient is an interface for making HTTP requests, used for dependency injection.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var httpClient HTTPClient = &http.Client{} // Default to the real HTTP client

// SetHTTPClient sets the global HTTP client for use in handlers.
func SetHTTPClient(client HTTPClient) {
	httpClient = client
}

// VehiclePositionsSimple fetches and displays simple vehicle positions from the API.
func (m *Repository) VehiclePositionsSimple(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	urlBase := apiBaseURL + "/colectivos/vehiclePositionsSimple"

	req, err := http.NewRequest("GET", urlBase, nil)
	if err != nil {
		stringMap["error"] = err.Error()
	}

	q := req.URL.Query()
	q.Add("client_id", clientID)
	q.Add("client_secret", clientSecret)

	req.URL.RawQuery = q.Encode()

	resp, err := httpClient.Do(req) // Use the injected HTTP client
	if err != nil {
		stringMap["error"] = err.Error()
	}

	if resp.StatusCode != http.StatusOK {
		stringMap["error"] = fmt.Sprintf("Error getting request, response code: %d", resp.StatusCode)
	}

	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			log.Println("error closing response body:", cerr)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		stringMap["error"] = err.Error()
	}

	stringMap["response"] = string(body)

	err = render.RenderTemplate(w, r, "positionsimple.page.tmpl", &model.TemplateData{
		StringMap: stringMap,
	})
	if err != nil {
		helpers.ServerError(w, err)
	}
}

// SearchLine renders the search page for bus lines.
func (m *Repository) SearchLine(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	data["line"] = ""

	err := render.RenderTemplate(w, r, "search.page.tmpl", &model.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
	if err != nil {
		helpers.ServerError(w, err)
	}
}

// TODO: consider to use a json for the form values
// by doing this we could expose the SearchLine as an API request
// PostSearchLine handles POST requests for searching bus lines.
func (m *Repository) PostSearchLine(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	line := r.Form.Get("line")

	form := forms.New(r.PostForm)
	form.Required("line")

	if !form.Valid() {
		err := render.RenderTemplate(w, r, "search.page.tmpl", &model.TemplateData{
			Form: form,
			Data: data,
		})
		if err != nil {
			helpers.ServerError(w, err)
		}

		return
	}

	result := utils.SearchLine(line, m.App.DataCache.Routes)
	resultString := strings.ReplaceAll(fmt.Sprintf("%+v", result), "} {", "} <br> {")

	data["result"] = template.HTML(resultString)
	data["line"] = line

	err = render.RenderTemplate(w, r, "search.page.tmpl", &model.TemplateData{
		Form: form,
		Data: data,
	})
	if err != nil {
		helpers.ServerError(w, err)
	}
}

// AllowedParking renders the allowed parking page.
func (m *Repository) AllowedParking(w http.ResponseWriter, r *http.Request) {
	googleMapsAPIKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	data := map[string]interface{}{
		"google_maps_api_key": googleMapsAPIKey,
	}
	err := render.RenderTemplate(w, r, "allowedparking.page.tmpl", &model.TemplateData{
		Data: data,
	})
	if err != nil {
		helpers.ServerError(w, err)
	}
}

// PostAllowedParking handles POST requests for allowed parking queries.
// It validates input, calls the ParkingRules API, and renders the result.
func (m *Repository) PostAllowedParking(w http.ResponseWriter, r *http.Request) {
	// Check for nil repository or config
	if m == nil || m.App == nil {
		log.Println("Repository or AppConfig is nil in PostAllowedParking")
		helpers.ServerError(w, fmt.Errorf("repository or AppConfig is nil"))
		return
	}

	data := make(map[string]interface{})

	// Parse the form data
	if err := r.ParseForm(); err != nil {
		log.Println("Error parsing form data:", err)
		helpers.ServerError(w, err)
		return
	}

	// Validate latitude and longitude
	latStr := r.Form.Get("latitude")
	lonStr := r.Form.Get("longitude")
	if latStr == "" || lonStr == "" {
		log.Println("Missing latitude or longitude in form data")
		helpers.ClientError(w, http.StatusBadRequest)
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		log.Println("Invalid latitude value:", latStr)
		helpers.ClientError(w, http.StatusBadRequest)
		return
	}

	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		log.Println("Invalid longitude value:", lonStr)
		helpers.ClientError(w, http.StatusBadRequest)
		return
	}

	// Call the ParkingRules service
	rules, err := m.APIClient.ParkingRules(lat, lon)
	if err != nil && !errors.Is(err, services.ErrNoParkingRules) {
		log.Println("Error calling ParkingRules service:", err)
		helpers.ServerError(w, err)
		return
	}

	// Populate the data map for the template
	data["address"] = r.Form.Get("address")
	data["latitude"] = latStr
	data["longitude"] = lonStr
	if rules != nil {
		data["rules"] = rules
	}
	if errors.Is(err, services.ErrNoParkingRules) {
		data["error"] = "No parking rules found for the specified location."
	}

	googleMapsAPIKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	data["google_maps_api_key"] = googleMapsAPIKey

	// Render the allowed parking template with the data
	err = render.RenderTemplate(w, r, "allowedparking.page.tmpl", &model.TemplateData{
		Data: data,
	})
	if err != nil {
		log.Println("Error rendering template:", err)
		helpers.ServerError(w, err)
	}
}

// FIXME: this func doesn't work
// FeedGtfsFrequency fetches GTFS frequency data from the API and prints trip IDs.
func (m *Repository) FeedGtfsFrequency(w http.ResponseWriter, r *http.Request) {
	if m == nil {
		log.Println("Repository is nil in FeedGtfsFrequency")
		helpers.ServerError(w, fmt.Errorf("repository is nil"))
		return
	}

	stringMap := make(map[string]string)

	urlBase := apiBaseURL + "/colectivos/feed-gtfs-frequency"

	req, err := http.NewRequest("GET", urlBase, nil)
	if err != nil {
		stringMap["error"] = err.Error()
		return
	}

	q := req.URL.Query()
	q.Add("client_id", clientID)
	q.Add("client_secret", clientSecret)

	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		stringMap["error"] = err.Error()
	}

	if resp == nil {
		stringMap["error"] = "Response is nil"
		return
	}

	if resp.StatusCode != http.StatusOK {
		stringMap["error"] = fmt.Sprintf("Error getting request, response code: %d", resp.StatusCode)
	}

	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			log.Println("error closing response body:", cerr)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		stringMap["error"] = err.Error()
	}

	feed := gtfs.FeedMessage{}
	err = proto.Unmarshal(body, &feed)
	if err != nil {
		log.Fatal(err)
	}

	for _, entity := range feed.Entity {
		tripUpdate := entity.GetTripUpdate()
		trip := tripUpdate.GetTrip()
		tripId := trip.GetTripId()
		fmt.Printf("Trip ID: %s\n", tripId)
	}

	// stringMap["response"] = string(body)

	// render.RenderTemplate(w, "positionsimple.page.tmpl", &model.TemplateData{
	// 	StringMap: stringMap,
	// })
}
