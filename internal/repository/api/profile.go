package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/llchhh/spektr-account-api/domain"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// ProfileRepository provides methods for profile management.
type ProfileRepository struct {
	client  *http.Client
	baseURL string
}

// NewProfileRepository creates a new ProfileRepository instance.
func NewProfileRepository(baseURL string) *ProfileRepository {
	return &ProfileRepository{
		client:  &http.Client{},
		baseURL: baseURL,
	}
}

// sendRequest sends a GET request to the API with the provided parameters and decodes the response.
func (p *ProfileRepository) sendRequest(ctx context.Context, method string, arg1 interface{}) ([]byte, error) {
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
	requestURL := fmt.Sprintf("%s?%s", p.baseURL, params.Encode())
	log.Printf("Request URL: %s", requestURL)

	// Create the HTTP GET request
	req, err := http.NewRequestWithContext(ctx, "GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Send the request
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed request, status code: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}

// Profile fetches the profile data for a user.
func (p *ProfileRepository) Profile(ctx context.Context, suid string) (domain.Profile, error) {
	log.Printf("Fetching profile for user with suid: %s", suid)

	arg1 := struct {
		SUID string `json:"suid"`
	}{SUID: suid}

	body, err := p.sendRequest(ctx, "web_cabinet.get_user", arg1)
	if err != nil {
		return domain.Profile{}, err
	}

	// Define the API response structure
	var apiResponse struct {
		Error string `json:"error"`
		User  struct {
			Abonent struct {
				Name          string      `json:"name"`
				Tariff        string      `json:"__tarif"`
				Balance       string      `json:"__account"`
				MinimalPaySum json.Number `json:"minimal_pay_sum"`
				Email         string      `json:"email"`
				Phone         string      `json:"sms"`
				AllowInternet string      `json:"allow_internet"`
				ID            string      `json:"contract_number"`
			} `json:"abonent"`
		} `json:"user"`
	}

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return domain.Profile{}, fmt.Errorf("failed to parse API response: %w", err)
	}

	// Handle errors in the response
	if apiResponse.Error != "" {
		return domain.Profile{}, fmt.Errorf("API error: %s", apiResponse.Error)

	}

	// Map API response to domain.Profile
	minimalPaySum, _ := apiResponse.User.Abonent.MinimalPaySum.Float64()
	firstName, middleName, lastName := parseFullName(apiResponse.User.Abonent.Name)

	return domain.Profile{
		FirstName:      middleName,
		MiddleName:     lastName,
		LastName:       firstName,
		Tariff:         apiResponse.User.Abonent.Tariff,
		Balance:        parseBalance(apiResponse.User.Abonent.Balance),
		ToPay:          minimalPaySum,
		Email:          apiResponse.User.Abonent.Email,
		Phone:          apiResponse.User.Abonent.Phone,
		InternetStatus: parseInternetStatus(apiResponse.User.Abonent.AllowInternet),
		ID:             apiResponse.User.Abonent.ID,
	}, nil
}

// ChangePassword changes the password for a user.
func (p *ProfileRepository) ChangePassword(ctx context.Context, suid, newPassword string) error {
	log.Printf("Changing password for user with suid: %s", suid)

	arg1 := struct {
		SUID         string `json:"suid"`
		UserPassword string `json:"user_password"`
	}{SUID: suid, UserPassword: newPassword}

	body, err := p.sendRequest(ctx, "web_cabinet.set_user_info", arg1)
	if err != nil {
		return err
	}
	// Define the API response structure
	var apiResponse struct {
		Error string `json:"error"`
	}

	// If the response body is empty, we can handle it as a special case
	if len(body) == 0 {
		log.Printf("Empty response body received")
		return fmt.Errorf("received empty response body")
	}

	// Parse the API response
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		log.Printf("Failed to parse API response: %v", err)
		return fmt.Errorf("failed to parse API response: %w", err)
	}

	// Handle errors in the response
	if apiResponse.Error != "" {
		return fmt.Errorf("API error: %s", apiResponse.Error)
	}

	log.Printf("Password change successful for user with suid: %s", suid)
	return nil
}

// ChangePhone changes the phone number for a user.
func (p *ProfileRepository) ChangePhone(ctx context.Context, suid, newPhone string) error {
	log.Printf("Changing phone number for user with suid: %s", suid)

	arg1 := struct {
		SUID string `json:"suid"`
		SMS  string `json:"sms"`
	}{SUID: suid, SMS: newPhone}

	body, err := p.sendRequest(ctx, "web_cabinet.set_user_info", arg1)
	if err != nil {
		return err
	}
	// Define the API response structure
	var apiResponse struct {
		Error string `json:"error"`
	}

	// If the response body is empty, we can handle it as a special case
	if len(body) == 0 {
		log.Printf("Empty response body received")
		return fmt.Errorf("received empty response body")
	}

	// Parse the API response
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		log.Printf("Failed to parse API response: %v", err)
		return fmt.Errorf("failed to parse API response: %w", err)
	}

	// Handle errors in the response
	if apiResponse.Error != "" {
		return fmt.Errorf("API error: %s", apiResponse.Error)
	}

	log.Printf("Phone change successful for user with suid: %s", suid)
	return nil
}

// ChangeEmail changes the email address for a user.
func (p *ProfileRepository) ChangeEmail(ctx context.Context, suid, newEmail string) error {
	log.Printf("Changing email for user with suid: %s", suid)

	arg1 := struct {
		SUID  string `json:"suid"`
		Email string `json:"email"`
	}{SUID: suid, Email: newEmail}

	body, err := p.sendRequest(ctx, "web_cabinet.set_user_info", arg1)
	if err != nil {
		return err
	}
	// Define the API response structure
	var apiResponse struct {
		Error string `json:"error"`
	}

	// If the response body is empty, we can handle it as a special case
	if len(body) == 0 {
		log.Printf("Empty response body received")
		return fmt.Errorf("received empty response body")
	}

	// Parse the API response
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		log.Printf("Failed to parse API response: %v", err)
		return fmt.Errorf("failed to parse API response: %w", err)
	}

	// Handle errors in the response
	if apiResponse.Error != "" {
		return fmt.Errorf("API error: %s", apiResponse.Error)
	}

	log.Printf("Email change successful for user with suid: %s", suid)
	return nil
}

// Helper functions
func parseBalance(balanceStr string) float64 {
	re := regexp.MustCompile(`Баланс:\s*([\d.]+)`)
	matches := re.FindStringSubmatch(balanceStr)
	if len(matches) > 1 {
		balance, _ := strconv.ParseFloat(matches[1], 64)
		return balance
	}
	return 0.0
}

func parseInternetStatus(statusStr string) bool {
	return statusStr == "1"
}

func parseFullName(fullName string) (string, string, string) {
	nameParts := strings.Fields(fullName)
	switch len(nameParts) {
	case 1:
		return nameParts[0], "", ""
	case 2:
		return nameParts[0], "", nameParts[1]
	case 3:
		return nameParts[0], nameParts[1], nameParts[2]
	default:
		return nameParts[0], strings.Join(nameParts[1:len(nameParts)-1], " "), nameParts[len(nameParts)-1]
	}
}
