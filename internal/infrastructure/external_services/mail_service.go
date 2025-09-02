package external_services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/RealEskalate/G6-NewsBrief/internal/domain/contract"
)

// EmailService implements sending emails via SendGrid API
// API docs: https://docs.sendgrid.com/api-reference/mail-send/mail-send
// Minimal plain-text implementation.
type EmailService struct {
	APIKey   string
	From     string
	FromName string
	client   *http.Client
}

// NewEmailService creates a SendGrid email service client.
func NewEmailService(apiKey, from, fromName string) *EmailService {
	return &EmailService{
		APIKey:   apiKey,
		From:     from,
		FromName: fromName,
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

var _ contract.IEmailService = (*EmailService)(nil)

// sendgridRequest models the JSON payload for SendGrid mail send
type sendgridRequest struct {
	Personalizations []struct {
		To      []struct{ Email string `json:"email"` } `json:"to"`
		Subject string                               `json:"subject"`
	} `json:"personalizations"`
	From struct {
		Email string `json:"email"`
		Name  string `json:"name,omitempty"`
	} `json:"from"`
	Content []struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"content"`
}

func (es *EmailService) SendEmail(ctx context.Context, to, subject, body string) error {
	if es.APIKey == "" {
		return fmt.Errorf("sendgrid api key is not configured")
	}
	if es.From == "" {
		return fmt.Errorf("email from address is not configured")
	}

	payload := sendgridRequest{}
	p := struct {
		To      []struct{ Email string `json:"email"` } `json:"to"`
		Subject string                               `json:"subject"`
	}{}
	p.To = []struct{ Email string `json:"email"` }{{Email: to}}
	p.Subject = subject
	payload.Personalizations = []struct {
		To      []struct{ Email string `json:"email"` } `json:"to"`
		Subject string                               `json:"subject"`
	}{p}
	payload.From.Email = es.From
	payload.From.Name = es.FromName
	payload.Content = []struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	}{{Type: "text/plain", Value: body}}

	buf, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal sendgrid payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.sendgrid.com/v3/mail/send", bytes.NewReader(buf))
	if err != nil {
		return fmt.Errorf("failed to create sendgrid request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+es.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := es.client.Do(req)
	if err != nil {
		return fmt.Errorf("sendgrid request failed: %w", err)
	}
	defer resp.Body.Close()

	// SendGrid returns 202 Accepted on success
	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("sendgrid mail send failed: status=%d", resp.StatusCode)
	}
	return nil
}
