package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/llchhh/spektr-account-api/domain"
	"github.com/llchhh/spektr-account-api/internal/encryption"
	"github.com/llchhh/spektr-account-api/internal/rest/middleware"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

type AuthRepository struct {
	client  *http.Client
	baseURL string
}

// Response struct to match the response format that contains session_id
type AuthResponse struct {
	SessionID string `json:"session_id"`
}

func (a *AuthRepository) Login(ctx context.Context, user domain.Auth) (string, error) {
	// Marshal the user data into a JSON object for arg1
	if middleware.IsSuspicious(user.Login) {
		return "", fmt.Errorf("invalid input detected")
	}

	if middleware.IsSuspicious(user.Password) {
		return "", fmt.Errorf("invalid input detected")
	}
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

	// Ensure the session_id is present in the response
	if result.SessionID == "" {
		return "", errors.New("session_id not found in the response")
	}
	encryptedToken, err := encryption.EncryptToken(os.Getenv("ENCRYPTION_KEY"), result.SessionID)
	// Return the session ID
	decrypted, _ := encryption.DecryptToken(os.Getenv("ENCRYPTION_KEY"), encryptedToken)
	fmt.Println(decrypted)
	return encryptedToken, nil
}

func NewAuthRepository(baseURL string) *AuthRepository {
	return &AuthRepository{
		client:  &http.Client{},
		baseURL: baseURL,
	}
}
