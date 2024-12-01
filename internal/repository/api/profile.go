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

type ProfileRepository struct {
	client  *http.Client
	baseURL string
}

// Fetches the profile data using the given token.
func (p *ProfileRepository) Profile(ctx context.Context, suid string) (domain.Profile, error) {
	// Log the incoming request for profiling
	log.Printf("Fetching profile for user with suid: %s", suid)

	// Construct the `arg1` parameter as a JSON string
	arg1 := fmt.Sprintf(`{"suid":"%s"}`, suid)

	// Construct query parameters
	params := url.Values{}
	params.Add("format", "json")
	params.Add("context", "web")
	params.Add("model", "users")
	params.Add("method1", "web_cabinet.get_user")
	params.Add("arg1", arg1)

	// Construct the request URL
	requestURL := fmt.Sprintf("%s?%s", p.baseURL, params.Encode())
	log.Printf("Request URL: %s", requestURL)

	// Create the HTTP GET request
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return domain.Profile{}, fmt.Errorf("failed to create request: %w", err)
	}

	// Set context and send the request
	req = req.WithContext(ctx)
	resp, err := p.client.Do(req)
	if err != nil {
		log.Printf("Request error: %v", err)
		return domain.Profile{}, fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to fetch profile, status code: %d", resp.StatusCode)
		return domain.Profile{}, fmt.Errorf("failed to fetch profile, status code: %d", resp.StatusCode)
	}

	// Parse the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return domain.Profile{}, fmt.Errorf("failed to read response body: %w", err)
	}

	// Define a struct to map the API response
	var apiResponse struct {
		Error      string `json:"error"`
		LoginError string `json:"login_error"`
		User       struct {
			Abonent struct {
				Name          string      `json:"name"`
				Tariff        string      `json:"__tarif"`
				Balance       string      `json:"__account"`
				MinimalPaySum json.Number `json:"minimal_pay_sum"` // Use json.Number for flexible handling
				Email         string      `json:"email"`
				Phone         string      `json:"sms"`
				AllowInternet string      `json:"allow_internet"`
				ID            string      `json:"contract_number"`
			} `json:"abonent"`
		} `json:"user"`
	}

	// Unmarshal the response into the local variable
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		log.Printf("Error unmarshalling response: %v", err)
		return domain.Profile{}, fmt.Errorf("failed to unmarshal API response: %w", err)
	}

	// If there's an error message in the response, handle it
	if apiResponse.Error != "" {
		log.Printf("API error: %s", apiResponse.Error)
		if apiResponse.Error == "Необходимо авторизоваться" {
			// Handle invalid token (authentication required)
			log.Println("Invalid token provided.")
			return domain.Profile{}, domain.ErrSessionExpired
		}
		return domain.Profile{}, fmt.Errorf("API error: %s", apiResponse.Error)
	}

	// Parse MinimalPaySum to float64
	minimalPaySum, err := apiResponse.User.Abonent.MinimalPaySum.Float64()
	if err != nil {
		log.Printf("Error parsing minimal_pay_sum: %v", err)
		minimalPaySum = 0.0 // Default value if parsing fails
	}

	// Parse the full name into firstName, middleName, and lastName
	firstName, middleName, lastName := parseFullName(apiResponse.User.Abonent.Name)

	// Map API response to the domain.Profile struct
	profile := domain.Profile{
		FirstName:      middleName,
		MiddleName:     lastName,
		LastName:       firstName,
		Tariff:         apiResponse.User.Abonent.Tariff,
		Balance:        parseBalance(apiResponse.User.Abonent.Balance), // Custom function to parse balance string
		ToPay:          minimalPaySum,
		Email:          apiResponse.User.Abonent.Email,
		Phone:          apiResponse.User.Abonent.Phone,
		InternetStatus: parseInternetStatus(apiResponse.User.Abonent.AllowInternet), // Convert '1'/'0' to true/false
		ID:             apiResponse.User.Abonent.ID,
	}

	log.Printf("Profile fetched successfully: %s %s %s", firstName, middleName, lastName)
	return profile, nil
}

// Helper function to parse the balance string (example: "№ 10001466 Баланс: 390.67")
func parseBalance(balanceStr string) float64 {
	// Extract the balance part from the string
	re := regexp.MustCompile(`Баланс:\s*([\d.]+)`)
	matches := re.FindStringSubmatch(balanceStr)
	if len(matches) > 1 {
		// Convert the balance value to float64
		balance, err := strconv.ParseFloat(matches[1], 64)
		if err == nil {
			log.Printf("Parsed balance: %f", balance)
			return balance
		}
	}
	log.Println("Failed to parse balance.")
	return 0.0
}

// Helper function to parse internet status
func parseInternetStatus(statusStr string) bool {
	if statusStr == "1" {
		log.Println("Internet access allowed.")
		return true
	}
	log.Println("Internet access not allowed.")
	return false
}

// Helper function to parse full name into first name, middle name, and last name
func parseFullName(fullName string) (string, string, string) {
	// Split the full name by spaces
	nameParts := strings.Fields(fullName)

	// Handle cases based on the number of name parts
	switch len(nameParts) {
	case 1:
		// Only last name provided
		return nameParts[0], "", ""
	case 2:
		// Last name and first name
		return nameParts[0], "", nameParts[1]
	case 3:
		// Last name, middle name, and first name
		return nameParts[0], nameParts[1], nameParts[2]
	default:
		// Last name, middle name(s), and first name
		return nameParts[0], strings.Join(nameParts[1:len(nameParts)-1], " "), nameParts[len(nameParts)-1]
	}
}

// ChangePassword отправляет запрос на изменение пароля.
func (p *ProfileRepository) ChangePassword(ctx context.Context, suid string, newPassword string) error {
	// Log the incoming request for password change
	log.Printf("Requesting password change for user with suid: %s", suid)

	// Validate input parameters
	if suid == "" {
		log.Println("ChangePassword failed: missing suid")
		return fmt.Errorf("suid is required")
	}
	// Construct the `arg1` parameter as a JSON string
	arg1 := fmt.Sprintf(`{"suid":"%s", "user_password":"%s"}`, suid, newPassword)

	// Construct query parameters
	params := url.Values{}
	params.Add("format", "json")
	params.Add("context", "web")
	params.Add("model", "users")
	params.Add("method1", "web_cabinet.set_user_info")
	params.Add("arg1", arg1)

	// Construct the request URL
	requestURL := fmt.Sprintf("%s?%s", p.baseURL, params.Encode())
	log.Printf("Request URL: %s", requestURL)

	// Create the HTTP GET request
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set context and send the request
	req = req.WithContext(ctx)
	resp, err := p.client.Do(req)
	if err != nil {
		log.Printf("Request error: %v", err)
		return fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to change password, status code: %d", resp.StatusCode)
		return fmt.Errorf("failed to change password, status code: %d", resp.StatusCode)
	}

	// Parse the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Define a struct to map the API response
	var apiResponse struct {
		Error string `json:"error"`
	}

	// Unmarshal the response into the struct
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		log.Printf("Error unmarshalling response: %v", err)
		return fmt.Errorf("failed to unmarshal API response: %w", err)
	}

	// Handle API errors
	if apiResponse.Error != "" {
		log.Printf("API error: %s", apiResponse.Error)
		if apiResponse.Error == "Необходимо авторизоваться" {
			// Handle invalid token (authentication required)
			log.Println("Invalid token provided.")
			return domain.ErrSessionExpired
		}
		return fmt.Errorf("API error: %s", apiResponse.Error)
	}

	log.Printf("Password successfully changed for user with suid: %s", suid)
	return nil
}

// ChangePhone sends a request to update the phone number for a user.
func (p *ProfileRepository) ChangePhone(ctx context.Context, suid string, newPhone string) error {
	// Log the incoming request for phone number change
	log.Printf("Requesting phone number change for user with suid: %s", suid)

	// Validate input parameters
	if suid == "" {
		log.Println("ChangePhone failed: missing suid")
		return fmt.Errorf("suid is required")
	}
	// Construct the `arg1` parameter as a JSON string
	arg1 := fmt.Sprintf(`{"suid":"%s", "sms":"%s"}`, suid, newPhone)

	// Construct query parameters
	params := url.Values{}
	params.Add("format", "json")
	params.Add("context", "web")
	params.Add("model", "users")
	params.Add("method1", "web_cabinet.set_user_info")
	params.Add("arg1", arg1)

	// Construct the request URL
	requestURL := fmt.Sprintf("%s?%s", p.baseURL, params.Encode())
	log.Printf("Request URL: %s", requestURL)

	// Create the HTTP GET request
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set context and send the request
	req = req.WithContext(ctx)
	resp, err := p.client.Do(req)
	if err != nil {
		log.Printf("Request error: %v", err)
		return fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to change phone number, status code: %d", resp.StatusCode)
		return fmt.Errorf("failed to change phone number, status code: %d", resp.StatusCode)
	}

	// Parse the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Define a struct to map the API response
	var apiResponse struct {
		Error string `json:"error"`
	}

	// Unmarshal the response into the struct
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		log.Printf("Error unmarshalling response: %v", err)
		return fmt.Errorf("failed to unmarshal API response: %w", err)
	}

	// Handle API errors
	if apiResponse.Error != "" {
		log.Printf("API error: %s", apiResponse.Error)
		if apiResponse.Error == "Необходимо авторизоваться" {
			// Handle invalid token (authentication required)
			log.Println("Invalid token provided.")
			return domain.ErrSessionExpired
		}
		return fmt.Errorf("API error: %s", apiResponse.Error)
	}

	log.Printf("Phone number successfully changed for user with suid: %s", suid)
	return nil
}

// ChangeEmail sends a request to update the email address for a user.
func (p *ProfileRepository) ChangeEmail(ctx context.Context, suid string, newEmail string) error {
	// Log the incoming request for email address change
	log.Printf("Requesting email address change for user with suid: %s", suid)

	// Validate input parameters
	if suid == "" {
		log.Println("ChangeEmail failed: missing suid")
		return fmt.Errorf("suid is required")
	}
	// Construct the `arg1` parameter as a JSON string
	arg1 := fmt.Sprintf(`{"suid":"%s", "email":"%s"}`, suid, newEmail)

	// Construct query parameters
	params := url.Values{}
	params.Add("format", "json")
	params.Add("context", "web")
	params.Add("model", "users")
	params.Add("method1", "web_cabinet.set_user_info")
	params.Add("arg1", arg1)

	// Construct the request URL
	requestURL := fmt.Sprintf("%s?%s", p.baseURL, params.Encode())
	log.Printf("Request URL: %s", requestURL)

	// Create the HTTP GET request
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set context and send the request
	req = req.WithContext(ctx)
	resp, err := p.client.Do(req)
	if err != nil {
		log.Printf("Request error: %v", err)
		return fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to change email address, status code: %d", resp.StatusCode)
		return fmt.Errorf("failed to change email address, status code: %d", resp.StatusCode)
	}

	// Parse the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return fmt.Errorf("failed to read response body: %w", err)
	}
	fmt.Println(suid, newEmail)
	fmt.Println(string(body))
	// Define a struct to map the API response
	var apiResponse struct {
		Error string `json:"error"`
	}

	// Unmarshal the response into the struct
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		log.Printf("Error unmarshalling response: %v", err)
		return fmt.Errorf("failed to unmarshal API response: %w", err)
	}

	// Handle API errors
	if apiResponse.Error != "" {
		log.Printf("API error: %s", apiResponse.Error)
		if apiResponse.Error == "Необходимо авторизоваться" {
			// Handle invalid token (authentication required)
			log.Println("Invalid token provided.")
			return domain.ErrSessionExpired
		}
		return fmt.Errorf("API error: %s", apiResponse.Error)
	}

	log.Printf("Email address successfully changed for user with suid: %s", suid)
	return nil
}

// NewProfileRepository creates a new ProfileRepository instance.
func NewProfileRepository(baseURL string) *ProfileRepository {
	return &ProfileRepository{
		client:  &http.Client{},
		baseURL: baseURL,
	}
}
