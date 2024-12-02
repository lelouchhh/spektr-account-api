package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/llchhh/spektr-account-api/domain"
	"io"
	"log"
	"net/http"
	"net/url"
)

// Notification represents the API notification structure.
type Notification struct {
	Text string `json:"text"`
	Type string `json:"type"`
}

// APIResponse represents a generic response structure from the API.
type APIResponse struct {
	Notifications []Notification `json:"notifications,omitempty"`
	Error         string         `json:"error,omitempty"`
}

// NotificationRepository handles fetching notifications from the API.
type NotificationRepository struct {
	client  *http.Client
	baseURL string
}

// GetNotifications fetches the notifications for the specified user.
func (n *NotificationRepository) GetNotifications(ctx context.Context, suid string) ([]domain.Notification, error) {
	log.Printf("Fetching notifications for user with suid: %s", suid)

	// Construct the payload as JSON
	payload := map[string]string{
		"suid": suid,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal payload: %v", err)
		return nil, fmt.Errorf("failed to prepare request payload: %w", err)
	}

	// Construct query parameters
	params := url.Values{
		"format":  {"json"},
		"context": {"web"},
		"model":   {"users"},
		"method1": {"web_cabinet.get_notifications_for_user"},
		"arg1":    {string(jsonPayload)},
	}

	// Build the request URL
	requestURL := fmt.Sprintf("%s?%s", n.baseURL, params.Encode())
	log.Printf("Request URL: %s", requestURL)

	// Create the HTTP GET request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Execute the HTTP request
	resp, err := n.client.Do(req)
	if err != nil {
		log.Printf("Request error: %v", err)
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Validate response status
	if resp.StatusCode != http.StatusOK {
		log.Printf("Non-OK HTTP status: %d", resp.StatusCode)
		return nil, fmt.Errorf("unexpected response status: %d", resp.StatusCode)
	}

	// Parse the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Use json.RawMessage to handle dynamic response
	var rawMessage json.RawMessage
	if err := json.Unmarshal(body, &rawMessage); err != nil {
		log.Printf("Failed to parse raw response: %v", err)
		return nil, fmt.Errorf("failed to parse API response: %w", err)
	}

	// Try to unmarshal into APIResponse first
	var apiResponse APIResponse
	if err := json.Unmarshal(rawMessage, &apiResponse); err == nil {
		// Handle error field
		if apiResponse.Error != "" {
			log.Printf("API returned an error: %s", apiResponse.Error)
			switch apiResponse.Error {
			case "Необходимо авторизоваться":
				return nil, domain.ErrSessionExpired
			default:
				return nil, errors.New(apiResponse.Error)
			}
		}
		// Convert API notifications to domain notifications
		notifications := make([]domain.Notification, len(apiResponse.Notifications))
		for i, apiNotification := range apiResponse.Notifications {
			notifications[i] = domain.Notification{
				Body: apiNotification.Text,
				Type: apiNotification.Type,
			}
		}
		log.Printf("Fetched %d notifications successfully", len(notifications))
		return notifications, nil
	}

	// If APIResponse fails, try to parse as an array
	var notificationArray []Notification
	if err := json.Unmarshal(rawMessage, &notificationArray); err == nil {
		// Convert API notifications to domain notifications
		notifications := make([]domain.Notification, len(notificationArray))
		for i, apiNotification := range notificationArray {
			notifications[i] = domain.Notification{
				Body: apiNotification.Text,
				Type: apiNotification.Type,
			}
		}
		log.Printf("Fetched %d notifications successfully from array response", len(notifications))
		return notifications, nil
	}

	// If all attempts fail, return an error
	log.Printf("Failed to parse API response: %v", err)
	return nil, fmt.Errorf("failed to parse API response: %w", err)
}

// NewNotificationRepository creates a new instance of NotificationRepository.
func NewNotificationRepository(baseURL string) *NotificationRepository {
	return &NotificationRepository{
		client:  &http.Client{},
		baseURL: baseURL,
	}
}
