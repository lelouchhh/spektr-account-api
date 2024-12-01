package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/llchhh/spektr-account-api/domain"
	"io/ioutil"
	"net/http"
	"net/url"
)

type AuthRepository struct {
	client  *http.Client
	baseURL string
}

func (a *AuthRepository) RequestPasswordResetToken(ctx context.Context, login string) error {
	arg1 := fmt.Sprintf(`{"login":"%s", "base_url":"null"}`, login)

	// Construct the query parameters in the exact order required
	params := url.Values{}
	params.Add("format", "json")
	params.Add("context", "web")
	params.Add("model", "users")
	params.Add("method1", "web_cabinet.reset_password")
	params.Add("arg1", arg1) // Add the JSON string for arg1

	// Construct the full URL with the query parameters
	requestURL := fmt.Sprintf("%s?%s", a.baseURL, params.Encode())

	// Create a new HTTP GET request (even though you're sending data in the query, it's a GET request)
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set the context for the request
	req = req.WithContext(ctx)

	// Send the request using http.Client
	resp, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()

	// Check if the response status is 200 (OK)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to authenticate, status code: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse the response body into a predefined struct
	var result AuthResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}
	if result.Error != "" {
		return domain.ErrInvalidCredentials
	}
	return nil
}

// UpdatePassword updates the user's password using the provided reset token, user ID, and new password.
func (a *AuthRepository) UpdatePassword(ctx context.Context, token, password string) error {
	// Create the payload with the token, user ID, and the new password
	arg1 := fmt.Sprintf(`{"token":"%s", "uid":"324", "psw1":"%s", "psw2":"%s"}`, token, password, password)

	// Construct the query parameters
	params := url.Values{}
	params.Add("format", "json")
	params.Add("context", "web")
	params.Add("model", "users")
	params.Add("method1", "web_cabinet.submit_password")
	params.Add("arg1", arg1) // Add the JSON string for arg1

	// Construct the full URL with the query parameters
	requestURL := fmt.Sprintf("%s?%s", a.baseURL, params.Encode())

	// Create a new HTTP GET request (even though you're sending data in the query, it's a GET request)
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set the context for the request
	req = req.WithContext(ctx)

	// Send the request using http.Client
	resp, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()

	// Check if the response status is 200 (OK)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update password, status code: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse the response body into a predefined struct
	var result AuthResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Handle errors from the API response
	if result.Error != "" {
		return domain.ErrInvalidCredentials
	}

	return nil
}

type AuthResponse struct {
	Error     string `json:"error"`
	SessionID string `json:"session_id"`
}

func (a *AuthRepository) Login(ctx context.Context, user domain.Auth) (string, error) {
	arg1JSON, err := json.Marshal(user)
	if err != nil {
		return "", fmt.Errorf("failed to marshal user data: %w", err)
	}

	// Convert the JSON object into a string to use in the query parameter
	arg1 := string(arg1JSON)

	// Construct the query parameters in the exact order required
	params := url.Values{}
	params.Add("format", "json")
	params.Add("context", "web")
	params.Add("model", "users")
	params.Add("method1", "web_cabinet.login")
	params.Add("arg1", arg1) // Add the JSON string for arg1

	// Construct the full URL with the query parameters
	requestURL := fmt.Sprintf("%s?%s", a.baseURL, params.Encode())

	// Create a new HTTP GET request (even though you're sending data in the query, it's a GET request)
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set the context for the request
	req = req.WithContext(ctx)

	// Send the request using http.Client
	resp, err := a.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()

	// Check if the response status is 200 (OK)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to authenticate, status code: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse the response body into a predefined struct
	var result AuthResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}
	if result.Error != "" {
		return "", domain.ErrInvalidCredentials
	}
	// Ensure the session_id is present in the response
	if result.SessionID == "" {
		return "", errors.New("session_id not found in the response")
	}
	return result.SessionID, nil
}

func NewAuthRepository(baseURL string) *AuthRepository {
	return &AuthRepository{
		client:  &http.Client{},
		baseURL: baseURL,
	}
}
