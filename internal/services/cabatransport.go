// Package services provides API clients and business logic for external services used by Bamos.
package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type (
	// *** CABA API transport client definition ***

	// APIClient defines the interface for interacting with the CABA transport API.
	APIClient interface {
		// ParkingRules fetches parking rules for a given latitude and longitude.
		ParkingRules(lat, long float64) (SimplifiedRules, error)
	}

	// Client implements APIClient and provides methods to interact with the CABA transport API.
	Client struct {
		ClientID     string
		ClientSecret string
		BaseURL      string
		HTTPClient   HTTPClient
	}

	// HTTPClient is an interface for making HTTP requests, used for dependency injection.
	HTTPClient interface {
		Do(req *http.Request) (*http.Response, error)
	}

	// *** Data types ***

	// ParkingRulesResponse represents the response from the parking rules API
	ParkingRulesResponse struct {
		TotalFull int        `json:"totalFull"`
		Instances []Instance `json:"instancias"`
		Total     int        `json:"total"`
	}

	// Instance represents a single parking rule instance
	Instance struct {
		Name     string `json:"nombre"`
		ClassID  string `json:"claseId"`
		Class    string `json:"clase"`
		ID       string `json:"id"`
		Distance string `json:"distancia"`
		Rules    Rule   `json:"contenido"`
	}

	// Rule represents the rules associated with a parking instance
	Rule struct {
		Detail []RuleDetail `json:"contenido"`
	}

	// RuleDetail represents the details of a parking rule
	RuleDetail struct {
		NameID    string `json:"nombreId"`
		Name      string `json:"nombre"`
		Possition string `json:"posicion"`
		Value     string `json:"valor"`
	}

	// SimplifiedRules represents a simplified version of parking rules
	SimplifiedRules map[string][]string
)

// ErrNoParkingRules is returned when no parking rules are found for the given coordinates.
var ErrNoParkingRules = errors.New("no parking rules found")

const (
	// BaseURL is the base URL for the API
	BaseURL = "https://apitransporte.buenosaires.gob.ar"
	// DefaultRetries is the default number of retries for API requests
	DefaultRetries = 3
	// DefaultTimeout is the default timeout for API requests
	DefaultTimeout = 3 * time.Second
)

// NewAPIClient creates and returns a new Client for the CABA transport API.
func NewAPIClient() *Client {
	clientID := os.Getenv("CABA_CLIENT_ID")
	clientSecret := os.Getenv("CABA_CLIENT_SECRET")
	return &Client{
		BaseURL:      BaseURL,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		HTTPClient: &http.Client{
			Timeout: DefaultTimeout,
		},
	}
}

// ParkingRules fetches parking rules for the specified latitude and longitude from the CABA API.
// Returns a SimplifiedRules map or NoParkingRulesError if no rules are found.
func (c *Client) ParkingRules(lat, long float64) (SimplifiedRules, error) {
	response := ParkingRulesResponse{}
	path := c.BaseURL + "/transito/v1/estacionamientos"
	var resp *http.Response
	var err error

	params := url.Values{}
	params.Add("client_id", c.ClientID)
	params.Add("client_secret", c.ClientSecret)
	params.Add("x", fmt.Sprintf("%f", long))
	params.Add("y", fmt.Sprintf("%f", lat))
	params.Add("radio", "100")
	params.Add("formato", "json")
	params.Add("fullInfo", "true")

	// Retry logic
	for i := 0; i < DefaultRetries; i++ {
		req, err := http.NewRequest("GET", fmt.Sprintf("%s?%s", path, params.Encode()), nil)
		if err != nil {
			return nil, fmt.Errorf("error creating request: %w", err)
		}

		resp, err = c.HTTPClient.Do(req)
		if err == nil && resp != nil && resp.StatusCode == http.StatusOK {
			break
		}
	}

	if resp == nil || resp.Body == nil {
		return nil, ErrNoParkingRules
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	if len(response.Instances) == 0 {
		// Return a specific error if no parking rules are found
		return nil, ErrNoParkingRules
	}

	return simplifyRules(response.Instances), nil
}

// simplifyRules converts API response instances to a simplified rules map.
func simplifyRules(instances []Instance) SimplifiedRules {
	simplifiedRules := SimplifiedRules{}

	for _, instance := range instances {
		var address, rule string
		fields := make(map[string]string, 6)

		for _, ruleDetail := range instance.Rules.Detail {
			switch ruleDetail.NameID {
			case "calle", "altura", "permiso", "horario", "lado", "paridad":
				fields[ruleDetail.NameID] = ruleDetail.Value
			default:
				continue
			}
		}
		// Construct the address
		address = fmt.Sprintf("%s %s", fields["calle"], fields["altura"])
		if value, ok := fields["paridad"]; ok && value != "" {
			rule = fmt.Sprintf("Lado %s (%s): %s las %s.",
				fields["lado"], fields["paridad"], fields["permiso"], fields["horario"])
		} else {
			rule = fmt.Sprintf("Lado %s: %s las %s.",
				fields["lado"], fields["permiso"], fields["horario"])
		}

		rule = strings.ToUpper(rule[:1]) + strings.ToLower(rule[1:])

		// Add the rule to the map
		simplifiedRules[address] = append(simplifiedRules[address], rule)
	}

	return simplifiedRules
}
