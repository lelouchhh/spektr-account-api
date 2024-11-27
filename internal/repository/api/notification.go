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
)

// Define the Notification structure
type Notification struct {
	Text string `json:"text"`
	Type string `json:"type"`
}

type NotificationRepository struct {
	client  *http.Client
	baseURL string
}

// Fetches the notifications for the given user.
func (n *NotificationRepository) GetNotifications(ctx context.Context, suid string) ([]domain.Notification, error) {
	// Log the incoming request for fetching notifications
	log.Printf("Fetching notifications for user with suid: %s", suid)

	// Construct the `arg1` parameter as a JSON string
	arg1 := fmt.Sprintf(`{"suid":"%s"}`, suid)

	// Construct query parameters
	params := url.Values{}
	params.Add("format", "json")
	params.Add("context", "web")
	params.Add("model", "users")
	params.Add("method1", "web_cabinet.get_notifications_for_user")
	params.Add("arg1", arg1)

	// Construct the request URL
	requestURL := fmt.Sprintf("%s?%s", n.baseURL, params.Encode())
	log.Printf("Request URL: %s", requestURL)

	// Create the HTTP GET request
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set context and send the request
	req = req.WithContext(ctx)
	resp, err := n.client.Do(req)
	if err != nil {
		log.Printf("Request error: %v", err)
		return nil, fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to fetch notifications, status code: %d", resp.StatusCode)
		return nil, domain.ErrInternalServerError
	}

	// Parse the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Define the response structure
	var apiResponse struct {
		Notifications []Notification
	}
	var errorResponse struct {
		Error string `json:"error"`
	}
	if err := json.Unmarshal(body, &errorResponse); err != nil {
		log.Printf("Error unmarshalling response: %v", err)
	}
	// Handle errors if any
	if errorResponse.Error != "" {
		log.Printf("API error: %s", errorResponse.Error)
		return nil, domain.ErrSessionExpired
	}
	if err := json.Unmarshal(body, &apiResponse.Notifications); err != nil {
		log.Printf("Error unmarshalling response: %v", err)
		return nil, fmt.Errorf("failed to unmarshal API response: %w", err)
	}

	// Map the notifications to the domain.Notification struct
	var notifications []domain.Notification
	for _, apiNotification := range apiResponse.Notifications {
		notifications = append(notifications, domain.Notification{
			Body: apiNotification.Text, // Mapping 'text' to 'Body'
			Type: apiNotification.Type, // Mapping 'type' to 'Type'
		})
	}

	log.Printf("Fetched %d notifications successfully", len(notifications))
	return notifications, nil
}

// NewNotificationRepository creates a new NotificationRepository instance.
func NewNotificationRepository(baseURL string) *NotificationRepository {
	return &NotificationRepository{
		client:  &http.Client{},
		baseURL: baseURL,
	}
}
