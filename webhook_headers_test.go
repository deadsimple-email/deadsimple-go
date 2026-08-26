package deadsimple

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// capture records what the client actually put on the wire.
type capture struct {
	method string
	path   string
	body   map[string]any
}

// newTestClient serves `response` for one request and records it.
func newTestClient(t *testing.T, status int, response string, got *capture) *Client {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got.method = r.Method
		got.path = r.URL.Path
		raw, _ := io.ReadAll(r.Body)
		if len(raw) > 0 {
			_ = json.Unmarshal(raw, &got.body)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = io.WriteString(w, response)
	}))
	t.Cleanup(srv.Close)
	return New("dse_test", WithBaseURL(srv.URL))
}

func TestCreateWebhookSendsCustomHeaders(t *testing.T) {
	var got capture
	c := newTestClient(t, 201, `{"data":{"webhook_id":"wh_001","url":"https://example.com/hook",
		"events":["message.received"],"headers":{"Authorization":"***3456"}}}`, &got)

	wh, err := c.Webhooks.Create(context.Background(), &CreateWebhookParams{
		URL:     "https://example.com/hook",
		Events:  []string{"message.received"},
		Headers: map[string]string{"Authorization": "Bearer tok_abcdef123456"},
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	sent, ok := got.body["headers"].(map[string]any)
	if !ok || sent["Authorization"] != "Bearer tok_abcdef123456" {
		t.Fatalf("headers not sent, body was %v", got.body)
	}
	if wh.Headers["Authorization"] != "***3456" {
		t.Fatalf("masked headers not decoded, got %v", wh.Headers)
	}
}

func TestCreateWebhookOmitsEmptyHeaders(t *testing.T) {
	var got capture
	c := newTestClient(t, 201, `{"data":{"webhook_id":"wh_002","url":"https://example.com/hook"}}`, &got)

	if _, err := c.Webhooks.Create(context.Background(), &CreateWebhookParams{
		URL:    "https://example.com/hook",
		Events: []string{"message.received"},
	}); err != nil {
		t.Fatalf("create: %v", err)
	}

	if _, present := got.body["headers"]; present {
		t.Fatalf("headers key should be omitted, body was %v", got.body)
	}
}

func TestGetWebhookHeaders(t *testing.T) {
	var got capture
	c := newTestClient(t, 200,
		`{"data":{"webhook_id":"wh_001","headers":{"Authorization":"***3456"}}}`, &got)

	headers, err := c.Webhooks.GetHeaders(context.Background(), "wh_001")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.method != http.MethodGet || got.path != "/v1/webhooks/wh_001/headers" {
		t.Fatalf("unexpected request %s %s", got.method, got.path)
	}
	if headers["Authorization"] != "***3456" {
		t.Fatalf("got %v", headers)
	}
}

func TestSetWebhookHeaders(t *testing.T) {
	var got capture
	c := newTestClient(t, 200,
		`{"data":{"webhook_id":"wh_001","headers":{"Authorization":"***9999"}}}`, &got)

	headers, err := c.Webhooks.SetHeaders(context.Background(), "wh_001",
		map[string]string{"Authorization": "Bearer rotated_9999"})
	if err != nil {
		t.Fatalf("set: %v", err)
	}
	if got.method != http.MethodPut || got.path != "/v1/webhooks/wh_001/headers" {
		t.Fatalf("unexpected request %s %s", got.method, got.path)
	}
	sent, _ := got.body["headers"].(map[string]any)
	if sent["Authorization"] != "Bearer rotated_9999" {
		t.Fatalf("body was %v", got.body)
	}
	if headers["Authorization"] != "***9999" {
		t.Fatalf("got %v", headers)
	}
}

// A nil map must serialize as {} — the API requires the key, and "null" would
// be rejected rather than clearing the headers.
func TestSetWebhookHeadersNilClears(t *testing.T) {
	var got capture
	c := newTestClient(t, 200, `{"data":{"webhook_id":"wh_001","headers":{}}}`, &got)

	headers, err := c.Webhooks.SetHeaders(context.Background(), "wh_001", nil)
	if err != nil {
		t.Fatalf("clear: %v", err)
	}
	sent, ok := got.body["headers"].(map[string]any)
	if !ok || len(sent) != 0 {
		t.Fatalf("expected an empty object, body was %v", got.body)
	}
	if len(headers) != 0 {
		t.Fatalf("expected no headers back, got %v", headers)
	}
}
