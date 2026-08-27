// Package deadsimple provides a Go client for the Dead Simple Email API.
//
// Usage:
//
//	client := deadsimple.New("dse_your_api_key")
//
//	inbox, err := client.Inboxes.Create(context.Background(), &deadsimple.CreateInboxParams{
//	    DisplayName: "Support Bot",
//	})
//
//	_, err = client.Messages.Send(context.Background(), inbox.InboxID, &deadsimple.SendMessageParams{
//	    To:       []string{"user@example.com"},
//	    Subject:  "Hello",
//	    TextBody: "Sent by an AI agent.",
//	})
package deadsimple

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	DefaultBaseURL = "https://api.deadsimple.email"
	Version        = "0.3.0"
)

// Client is the Dead Simple Email API client.
type Client struct {
	apiKey  string
	baseURL string
	http    *http.Client

	Inboxes      *InboxService
	Messages     *MessageService
	Threads      *ThreadService
	Webhooks     *WebhookService
	Domains      *DomainService
	ApiKeys      *ApiKeyService
	Workspaces   *WorkspaceService
	Usage        *UsageService
	Drafts       *DraftService
	Templates    *TemplateService
	Filters      *FilterService
	Contacts     *ContactService
	Suppressions *SuppressionService
	Events       *EventService
	Attachments  *AttachmentService
}

// Option configures the client.
type Option func(*Client)

// WithBaseURL sets a custom base URL.
func WithBaseURL(url string) Option {
	return func(c *Client) { c.baseURL = url }
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) { c.http = hc }
}

// WithTimeout sets the HTTP timeout.
func WithTimeout(d time.Duration) Option {
	return func(c *Client) { c.http.Timeout = d }
}

// New creates a new Dead Simple Email API client.
func New(apiKey string, opts ...Option) *Client {
	c := &Client{
		apiKey:  apiKey,
		baseURL: DefaultBaseURL,
		http:    &http.Client{Timeout: 30 * time.Second},
	}
	for _, opt := range opts {
		opt(c)
	}
	c.Inboxes = &InboxService{c: c}
	c.Messages = &MessageService{c: c}
	c.Threads = &ThreadService{c: c}
	c.Webhooks = &WebhookService{c: c}
	c.Domains = &DomainService{c: c}
	c.ApiKeys = &ApiKeyService{c: c}
	c.Workspaces = &WorkspaceService{c: c}
	c.Usage = &UsageService{c: c}
	c.Drafts = &DraftService{c: c}
	c.Templates = &TemplateService{c: c}
	c.Filters = &FilterService{c: c}
	c.Contacts = &ContactService{c: c}
	c.Suppressions = &SuppressionService{c: c}
	c.Events = &EventService{c: c}
	c.Attachments = &AttachmentService{c: c}
	return c
}

// ── HTTP layer ──────────────────────────────────────────────────────────────

type apiResponse struct {
	Data  json.RawMessage `json:"data"`
	Error *apiError       `json:"error,omitempty"`
}

type apiError struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

func (c *Client) do(ctx context.Context, method, path string, body interface{}, params url.Values, idempotencyKey string) (json.RawMessage, error) {
	u := c.baseURL + path
	if len(params) > 0 {
		u += "?" + params.Encode()
	}

	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal body: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, u, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "deadsimple-go/"+Version)
	if idempotencyKey != "" {
		req.Header.Set("Idempotency-Key", idempotencyKey)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 204 {
		return nil, nil
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	var apiResp apiResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("parse response: %w (status %d)", err, resp.StatusCode)
	}

	if resp.StatusCode >= 400 {
		msg := "Unknown error"
		code := ""
		if apiResp.Error != nil {
			msg = apiResp.Error.Message
			code = apiResp.Error.Code
		}
		return nil, &Error{
			Message:    msg,
			Code:       code,
			StatusCode: resp.StatusCode,
		}
	}

	return apiResp.Data, nil
}

func decodeJSON[T any](data json.RawMessage) (*T, error) {
	var result T
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	return &result, nil
}
