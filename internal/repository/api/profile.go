package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/llchhh/spektr-account-api/domain"
	"github.com/llchhh/spektr-account-api/internal/encryption"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

type ProfileRepository struct {
	client  *http.Client
	baseURL string
}

// Fetches the profile data using the given token.
func (p *ProfileRepository) Profile(ctx context.Context, token string) (domain.Profile, error) {
	// Decrypt the token
	decryptedToken, err := encryption.DecryptToken(os.Getenv("ENCRYPTION_KEY"), token)
	if err != nil {
		return domain.Profile{}, fmt.Errorf("failed to decrypt token: %w", err)
	}

	// Construct query parameters
	params := url.Values{}
	params.Add("format", "json")
	params.Add("context", "web")
	params.Add("model", "users")
	params.Add("method1", "web_cabinet.get_profile")
	params.Add("arg1", decryptedToken)

	// Construct the request URL
	requestURL := fmt.Sprintf("%s?%s", p.baseURL, params.Encode())

	// Create the HTTP GET request
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return domain.Profile{}, fmt.Errorf("failed to create request: %w", err)
	}

	// Set context and send the request
	req = req.WithContext(ctx)
	resp, err := p.client.Do(req)
	if err != nil {
		return domain.Profile{}, fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		return domain.Profile{}, fmt.Errorf("failed to fetch profile, status code: %d", resp.StatusCode)
	}

	// Parse the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return domain.Profile{}, fmt.Errorf("failed to read response body: %w", err)
	}

	// Unmarshal the profile data
	var profile domain.Profile
	if err := json.Unmarshal(body, &profile); err != nil {
		return domain.Profile{}, fmt.Errorf("failed to unmarshal profile: %w", err)
	}

	return profile, nil
}

// Updates the user's password.
func (p *ProfileRepository) ChangePassword(ctx context.Context, token string, password string) error {
	decryptedToken, err := encryption.DecryptToken(os.Getenv("ENCRYPTION_KEY"), token)
	if err != nil {
		return fmt.Errorf("failed to decrypt token: %w", err)
	}

	// Construct query parameters
	params := url.Values{}
	params.Add("format", "json")
	params.Add("context", "web")
	params.Add("model", "users")
	params.Add("method1", "web_cabinet.change_password")
	params.Add("arg1", decryptedToken)
	params.Add("arg2", password)

	// Construct the request URL
	requestURL := fmt.Sprintf("%s?%s", p.baseURL, params.Encode())

	// Create the HTTP POST request
	req, err := http.NewRequest("POST", requestURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set context and send the request
	req = req.WithContext(ctx)
	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to change password, status code: %d", resp.StatusCode)
	}

	return nil
}

// Updates the user's email.
func (p *ProfileRepository) ChangeEmail(ctx context.Context, token string, email string) error {
	decryptedToken, err := encryption.DecryptToken(os.Getenv("ENCRYPTION_KEY"), token)
	if err != nil {
		return fmt.Errorf("failed to decrypt token: %w", err)
	}

	// Construct query parameters
	params := url.Values{}
	params.Add("format", "json")
	params.Add("context", "web")
	params.Add("model", "users")
	params.Add("method1", "web_cabinet.change_email")
	params.Add("arg1", decryptedToken)
	params.Add("arg2", email)

	// Construct the request URL
	requestURL := fmt.Sprintf("%s?%s", p.baseURL, params.Encode())

	// Create the HTTP POST request
	req, err := http.NewRequest("POST", requestURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set context and send the request
	req = req.WithContext(ctx)
	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to change email, status code: %d", resp.StatusCode)
	}

	return nil
}

// NewProfileRepository creates a new ProfileRepository instance.
func NewProfileRepository(baseURL string) *ProfileRepository {
	return &ProfileRepository{
		client:  &http.Client{},
		baseURL: baseURL,
	}
}
