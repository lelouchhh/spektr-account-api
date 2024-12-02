package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/llchhh/spektr-account-api/domain"
	"io"
	"log"
	"net/http"
	"net/url"
)

// APIResponse represents a generic response structure from the API.
type apiResponseRepair struct {
	Error string `json:"error,omitempty"`
}

// RepairRepository handles creating repairs through the API.
type RepairRepository struct {
	client  *http.Client
	baseURL string
}

func (r *RepairRepository) CreateRepair(ctx context.Context, token string, repair domain.Repair) error {
	log.Printf("Creating repair request for user with token: %s", token)

	// Construct the payload as a map
	payload := map[string]interface{}{
		"suid":        token,
		"ticket_text": repair.Text,
		"subj":        repair.Subject,
		"status":      "1", // Add a status field (1 for active)
	}

	// Use the sendRequest method to perform the request
	body, err := r.sendRequest(ctx, "web_cabinet.create_ticket", payload)
	if err != nil {
		log.Printf("Request failed: %v", err)
		return err
	}

	// Parse the response from the API
	var apiResponse apiResponseRepair
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		log.Printf("Failed to parse API response: %v", err)
		return fmt.Errorf("failed to parse API response: %w", err)
	}

	// Handle any errors from the API
	if apiResponse.Error != "" {
		log.Printf("API returned an error: %s", apiResponse.Error)
		return fmt.Errorf("API error: %s", apiResponse.Error)
	}

	log.Printf("Repair request created successfully")
	return nil
}

// sendRequest sends a GET request to the API with the specified method and payload.
func (r *RepairRepository) sendRequest(ctx context.Context, method string, arg1 interface{}) ([]byte, error) {
	// Serialize arg1 into JSON
	jsonData, err := json.Marshal(arg1)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize arg1 to JSON: %w", err)
	}

	// Construct query parameters
	params := url.Values{}
	params.Add("format", "json")
	params.Add("context", "web")
	params.Add("model", "users")
	params.Add("method1", method)
	params.Add("arg1", string(jsonData))

	// Build the request URL
	requestURL := fmt.Sprintf("%s?%s", r.baseURL, params.Encode())
	log.Printf("Request URL: %s", requestURL)

	// Create the HTTP GET request
	req, err := http.NewRequestWithContext(ctx, "GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Send the request
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed request, status code: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}

// NewRepairRepository creates a new instance of RepairRepository.
func NewRepairRepository(baseURL string) *RepairRepository {
	return &RepairRepository{
		client:  &http.Client{},
		baseURL: baseURL,
	}
}
