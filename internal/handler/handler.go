package handler

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"google.golang.org/protobuf/proto"

	"github.com/mayloo89/bamos/internal/config"
	"github.com/mayloo89/bamos/internal/driver"
	"github.com/mayloo89/bamos/internal/forms"
	"github.com/mayloo89/bamos/internal/helpers"
	"github.com/mayloo89/bamos/internal/model"
	"github.com/mayloo89/bamos/internal/render"
	"github.com/mayloo89/bamos/internal/repository"
	"github.com/mayloo89/bamos/internal/repository/database"
	"github.com/mayloo89/bamos/utils"
)

var Repo *Repository
var clientID = "REPLACED"
var clientSecret = "REPLACED"

type (
	Repository struct {
		App *config.AppConfig
		DB  repository.DBRepository
	}
)

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  database.NewPostgreDatabase(db.SQL, a),
	}
}

// NewHandler sets the repository for the handlers
func NewHandler(r *Repository) {
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {

	stringMap := make(map[string]string)
	stringMap["test"] = "Hello again."

	err := render.Template(w, r, "home.page.tmpl", &model.TemplateData{
		StringMap: stringMap,
	})
	if err != nil {
		helpers.ServerError(w, err)
	}
}

func (m *Repository) VehiclePositionsSimple(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	urlBase := "https://apitransporte.buenosaires.gob.ar/colectivos/vehiclePositionsSimple"

	req, err := http.NewRequest("GET", urlBase, nil)
	if err != nil {
		stringMap["error"] = err.Error()
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

	if resp.StatusCode != http.StatusOK {
		stringMap["error"] = fmt.Sprintf("Error getting request, response code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		stringMap["error"] = err.Error()
	}

	stringMap["response"] = string(body)

	err = render.Template(w, r, "positionsimple.page.tmpl", &model.TemplateData{
		StringMap: stringMap,
	})
	if err != nil {
		helpers.ServerError(w, err)
	}
}

func (m *Repository) SearchLine(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	data["line"] = ""

	err := render.Template(w, r, "search.page.tmpl", &model.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
	if err != nil {
		helpers.ServerError(w, err)
	}
}

// TODO: consider to use a json for the form values
// by doing this we could expose the SearchLine as an API request
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
		err := render.Template(w, r, "search.page.tmpl", &model.TemplateData{
			Form: form,
			Data: data,
		})
		if err != nil {
			helpers.ServerError(w, err)
		}

		return
	}

	result := utils.SearchLine(line, m.App.DataCache.Routes)
	resultString := strings.Replace(fmt.Sprintf("%+v", result), "} {", "} <br> {", -1)

	data["result"] = template.HTML(resultString)
	data["line"] = line

	err = render.Template(w, r, "search.page.tmpl", &model.TemplateData{
		Form: form,
		Data: data,
	})
	if err != nil {
		helpers.ServerError(w, err)
	}
}

// FIXME: this func doesn't work
func (m *Repository) FeedGtfsFrequency(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)

	urlBase := "https://apitransporte.buenosaires.gob.ar/colectivos/feed-gtfs-frequency"

	req, err := http.NewRequest("GET", urlBase, nil)
	if err != nil {
		stringMap["error"] = err.Error()
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

	if resp.StatusCode != http.StatusOK {
		stringMap["error"] = fmt.Sprintf("Error getting request, response code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

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

	// render.Template(w, "positionsimple.page.tmpl", &model.TemplateData{
	// 	StringMap: stringMap,
	// })
}
