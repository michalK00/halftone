package fcm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2/google"
	"net/http"
)

type Service struct {
	projectID string
	client    *http.Client
}

type SubscriptionRequest struct {
	Token          string            `json:"token"`
	UserID         string            `json:"userId"`
	UserAttributes map[string]string `json:"userAttributes,omitempty"`
}

type SendMessageRequest struct {
	UserIDs []string     `json:"userIds,omitempty"`
	Message *PushMessage `json:"message"`
}

type PushMessage struct {
	Title string            `json:"title"`
	Body  string            `json:"body"`
	Data  map[string]string `json:"data,omitempty"`
	URL   string            `json:"url,omitempty"`
}

type FCMMessage struct {
	Message FCMMessagePayload `json:"message"`
}

type FCMMessagePayload struct {
	Token        string            `json:"token,omitempty"`
	Notification FCMNotification   `json:"notification"`
	Data         map[string]string `json:"data,omitempty"`
	WebPush      *FCMWebPush       `json:"webpush,omitempty"`
}

type FCMNotification struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type FCMWebPush struct {
	Headers      map[string]string `json:"headers,omitempty"`
	Data         map[string]string `json:"data,omitempty"`
	Notification map[string]any    `json:"notification,omitempty"`
}

var userTokens = make(map[string]string)

func NewService(projectID string, credentialsJSON []byte) (*Service, error) {
	ctx := context.Background()

	config, err := google.JWTConfigFromJSON(
		credentialsJSON,
		"https://www.googleapis.com/auth/cloud-platform",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create JWT config: %w", err)
	}

	client := config.Client(ctx)

	return &Service{
		projectID: projectID,
		client:    client,
	}, nil
}

func (s *Service) Subscribe(req *SubscriptionRequest) error {
	userTokens[req.UserID] = req.Token

	fmt.Printf("Successfully registered user %s with token %s\n", req.UserID, req.Token[:20]+"...")
	return nil
}

func (s *Service) SendMessage(req *SendMessageRequest) error {
	var tokens []string

	if len(req.UserIDs) > 0 {
		for _, userID := range req.UserIDs {
			if token, exists := userTokens[userID]; exists {
				tokens = append(tokens, token)
			}
		}
	}

	if len(tokens) == 0 {
		return fmt.Errorf("no valid tokens found")
	}

	for _, token := range tokens {
		if err := s.sendToToken(token, req.Message); err != nil {
			fmt.Printf("Failed to send to token %s: %v\n", token[:20]+"...", err)
			continue
		}
		fmt.Printf("Successfully sent message to token %s\n", token[:20]+"...")
	}

	return nil
}

func (s *Service) sendToToken(token string, message *PushMessage) error {
	fcmMessage := FCMMessage{
		Message: FCMMessagePayload{
			Token: token,
			Notification: FCMNotification{
				Title: message.Title,
				Body:  message.Body,
			},
			Data: message.Data,
			WebPush: &FCMWebPush{
				Headers: map[string]string{
					"TTL": "86400", // 24 hours
				},
				Notification: map[string]any{
					"title": message.Title,
					"body":  message.Body,
					"icon":  "/icon-192x192.png",
				},
			},
		},
	}

	if message.URL != "" {
		if fcmMessage.Message.WebPush.Notification == nil {
			fcmMessage.Message.WebPush.Notification = make(map[string]any)
		}
		fcmMessage.Message.WebPush.Notification["click_action"] = message.URL

		if fcmMessage.Message.Data == nil {
			fcmMessage.Message.Data = make(map[string]string)
		}
		fcmMessage.Message.Data["url"] = message.URL
	}

	payload, err := json.Marshal(fcmMessage)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	url := fmt.Sprintf("https://fcm.googleapis.com/v1/projects/%s/messages:send", s.projectID)

	resp, err := s.client.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResponse map[string]any
		err := json.NewDecoder(resp.Body).Decode(&errorResponse)
		if err != nil {
			return err
		}
		return fmt.Errorf("FCM request failed with status %d: %v", resp.StatusCode, errorResponse)
	}

	return nil
}
