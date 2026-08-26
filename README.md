# Dead Simple Email — Go SDK

The official Go SDK for [Dead Simple Email](https://deadsimple.email), the email API for AI agents.

- **Typed** — structs for every resource, not `map[string]interface{}`
- **Context-aware** — every call takes a `context.Context`
- **Zero dependencies** — standard library only
- **Idempotency** — pass an idempotency key to create and send calls for safe retries
- **Webhook verification** — HMAC-SHA256 signature checking built in

## Install

```bash
go get github.com/deadsimple-email/deadsimple-go
```

## Quick start

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	deadsimple "github.com/deadsimple-email/deadsimple-go"
)

func main() {
	client := deadsimple.New(os.Getenv("DSE_API_KEY"))
	ctx := context.Background()

	inbox, err := client.Inboxes.Create(ctx, &deadsimple.CreateInboxParams{
		DisplayName: "Support Bot",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inbox:", inbox.Email)

	result, err := client.Messages.Send(ctx, inbox.InboxID, &deadsimple.SendMessageParams{
		To:       []string{"user@example.com"},
		Subject:  "Hello from my AI agent",
		TextBody: "This email was sent by an AI agent using Dead Simple Email.",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Sent:", result.MessageID)
}
```

## Services

`client.Inboxes`, `client.Messages`, `client.Threads`, `client.Drafts`, `client.Filters`,
`client.Templates`, `client.Webhooks`, `client.Domains`, `client.ApiKeys`,
`client.Workspaces`, `client.Contacts`, `client.Suppressions`, `client.Events`,
`client.Usage`.

## Verifying webhooks

Every delivery is signed with HMAC-SHA256 over `"<timestamp>.<raw_body>"` using the
per-webhook signing secret. Verify the raw body, before decoding it.

```go
func handler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if err := deadsimple.VerifyWebhookRequest(body, r.Header, signingSecret, 0); err != nil {
		http.Error(w, "bad signature", http.StatusUnauthorized)
		return
	}

	// Return 2xx fast, then do the work asynchronously.
	w.WriteHeader(http.StatusOK)
	go process(body)
}
```

`VerifyWebhookRequest` reads whichever signature header is present
(`X-DSE-Signature`, or `X-DSE-Webhook-Signature` with `X-DSE-Webhook-Timestamp`).
Pass `0` for the default 300-second replay tolerance.

Handlers must be idempotent: deliveries are retried on non-2xx or timeout, so the
same event can legitimately arrive more than once. Key on `message_id`.

## Custom delivery headers

If your endpoint enforces its own auth, register the headers it requires. They
are sent on every delivery attempt and every retry.

```go
wh, err := client.Webhooks.Create(ctx, &deadsimple.CreateWebhookParams{
	URL:     "https://your-app.com/webhook",
	Events:  []string{"message.received"},
	Headers: map[string]string{"Authorization": "Bearer your-endpoint-token"},
})

// Rotate later without recreating the webhook (which would change its secret).
_, err = client.Webhooks.SetHeaders(ctx, wh.WebhookID,
	map[string]string{"Authorization": "Bearer rotated"})
```

Values are write-only: `GetHeaders` and every other response returns them masked.
Up to 10 headers, 1024 characters per value, printable ASCII. `Content-Type`,
`Host`, `User-Agent`, hop-by-hop headers and `X-DSE-*` are reserved.

## Options

```go
client := deadsimple.New(apiKey,
	deadsimple.WithBaseURL("https://api.deadsimple.email"),
	deadsimple.WithHTTPClient(&http.Client{Timeout: 30 * time.Second}),
	deadsimple.WithTimeout(30*time.Second),
)
```

## Links

- [Documentation](https://deadsimple.email/docs.html)
- [API reference](https://deadsimple.email/api-reference.html)
- [Dashboard](https://app.deadsimple.email)

## License

MIT
