package deadsimple

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// ── Inboxes ─────────────────────────────────────────────────────────────────

type InboxService struct{ c *Client }

type CreateInboxParams struct {
	DisplayName    string   `json:"display_name,omitempty"`
	Domain         string   `json:"domain,omitempty"`
	LocalPart      string   `json:"local_part,omitempty"`
	Tags           []string `json:"tags,omitempty"`
	IdempotencyKey string   `json:"-"`
}

func (s *InboxService) Create(ctx context.Context, p *CreateInboxParams) (*Inbox, error) {
	key := ""
	if p != nil {
		key = p.IdempotencyKey
	}
	data, err := s.c.do(ctx, "POST", "/v1/inboxes", p, nil, key)
	if err != nil {
		return nil, err
	}
	return decodeJSON[Inbox](data)
}

type ListInboxesParams struct {
	Domain string
	Status string
	Tag    string
	Limit  int
	Cursor string
}

func (s *InboxService) List(ctx context.Context, p *ListInboxesParams) (*InboxList, error) {
	params := url.Values{}
	if p != nil {
		if p.Domain != "" {
			params.Set("domain", p.Domain)
		}
		if p.Status != "" {
			params.Set("status", p.Status)
		}
		if p.Tag != "" {
			params.Set("tag", p.Tag)
		}
		if p.Limit > 0 {
			params.Set("limit", strconv.Itoa(p.Limit))
		}
		if p.Cursor != "" {
			params.Set("cursor", p.Cursor)
		}
	}
	data, err := s.c.do(ctx, "GET", "/v1/inboxes", nil, params, "")
	if err != nil {
		return nil, err
	}
	return decodeJSON[InboxList](data)
}

func (s *InboxService) Get(ctx context.Context, inboxID string) (*Inbox, error) {
	data, err := s.c.do(ctx, "GET", fmt.Sprintf("/v1/inboxes/%s", inboxID), nil, nil, "")
	if err != nil {
		return nil, err
	}
	return decodeJSON[Inbox](data)
}

func (s *InboxService) Delete(ctx context.Context, inboxID string) error {
	_, err := s.c.do(ctx, "DELETE", fmt.Sprintf("/v1/inboxes/%s", inboxID), nil, nil, "")
	return err
}

func (s *InboxService) BulkCreate(ctx context.Context, inboxes []CreateInboxParams, idempotencyKey string) (*BulkCreateResult, error) {
	body := map[string]interface{}{"inboxes": inboxes}
	data, err := s.c.do(ctx, "POST", "/v1/inboxes/bulk", body, nil, idempotencyKey)
	if err != nil {
		return nil, err
	}
	return decodeJSON[BulkCreateResult](data)
}

// ── Messages ────────────────────────────────────────────────────────────────

type MessageService struct{ c *Client }

type SendMessageParams struct {
	To             []string          `json:"to"`
	Subject        string            `json:"subject"`
	TextBody       string            `json:"text_body,omitempty"`
	HTMLBody       string            `json:"html_body,omitempty"`
	CC             []string          `json:"cc,omitempty"`
	BCC            []string          `json:"bcc,omitempty"`
	ReplyTo        string            `json:"reply_to,omitempty"`
	Headers        map[string]string `json:"headers,omitempty"`
	SendAt         string            `json:"send_at,omitempty"`
	TemplateID     string            `json:"template_id,omitempty"`
	Variables      map[string]string `json:"variables,omitempty"`
	IdempotencyKey string            `json:"-"`
}

func (s *MessageService) Send(ctx context.Context, inboxID string, p *SendMessageParams) (*SendResult, error) {
	key := ""
	if p != nil {
		key = p.IdempotencyKey
	}
	data, err := s.c.do(ctx, "POST", fmt.Sprintf("/v1/inboxes/%s/messages", inboxID), p, nil, key)
	if err != nil {
		return nil, err
	}
	return decodeJSON[SendResult](data)
}

type ListMessagesParams struct {
	Q         string
	Direction string
	From      string
	Limit     int
	Cursor    string
}

func (s *MessageService) List(ctx context.Context, inboxID string, p *ListMessagesParams) (*MessageList, error) {
	params := url.Values{}
	if p != nil {
		if p.Q != "" {
			params.Set("q", p.Q)
		}
		if p.Direction != "" {
			params.Set("direction", p.Direction)
		}
		if p.From != "" {
			params.Set("from", p.From)
		}
		if p.Limit > 0 {
			params.Set("limit", strconv.Itoa(p.Limit))
		}
		if p.Cursor != "" {
			params.Set("cursor", p.Cursor)
		}
	}
	data, err := s.c.do(ctx, "GET", fmt.Sprintf("/v1/inboxes/%s/messages", inboxID), nil, params, "")
	if err != nil {
		return nil, err
	}
	return decodeJSON[MessageList](data)
}

func (s *MessageService) Get(ctx context.Context, inboxID, messageID string) (*Message, error) {
	data, err := s.c.do(ctx, "GET", fmt.Sprintf("/v1/inboxes/%s/messages/%s", inboxID, messageID), nil, nil, "")
	if err != nil {
		return nil, err
	}
	return decodeJSON[Message](data)
}

func (s *MessageService) Reply(ctx context.Context, inboxID, messageID string, p *ReplyParams) (*SendResult, error) {
	data, err := s.c.do(ctx, "POST", fmt.Sprintf("/v1/inboxes/%s/messages/%s/reply", inboxID, messageID), p, nil, "")
	if err != nil {
		return nil, err
	}
	return decodeJSON[SendResult](data)
}

func (s *MessageService) ReplyAll(ctx context.Context, inboxID, messageID string, p *ReplyParams) (*SendResult, error) {
	data, err := s.c.do(ctx, "POST", fmt.Sprintf("/v1/inboxes/%s/messages/%s/reply-all", inboxID, messageID), p, nil, "")
	if err != nil {
		return nil, err
	}
	return decodeJSON[SendResult](data)
}

type ReplyParams struct {
	TextBody string `json:"text_body,omitempty"`
	HTMLBody string `json:"html_body,omitempty"`
}

type ForwardParams struct {
	To       []string `json:"to"`
	TextBody string   `json:"text_body,omitempty"`
	HTMLBody string   `json:"html_body,omitempty"`
}

func (s *MessageService) Forward(ctx context.Context, inboxID, messageID string, p *ForwardParams) (*SendResult, error) {
	data, err := s.c.do(ctx, "POST", fmt.Sprintf("/v1/inboxes/%s/messages/%s/forward", inboxID, messageID), p, nil, "")
	if err != nil {
		return nil, err
	}
	return decodeJSON[SendResult](data)
}

// ── Threads ─────────────────────────────────────────────────────────────────

type ThreadService struct{ c *Client }

func (s *ThreadService) ListAll(ctx context.Context, inboxID, labels string, limit, offset int) (*ThreadList, error) {
	params := url.Values{"limit": {strconv.Itoa(limit)}, "offset": {strconv.Itoa(offset)}}
	if inboxID != "" {
		params.Set("inbox_id", inboxID)
	}
	if labels != "" {
		params.Set("labels", labels)
	}
	data, err := s.c.do(ctx, "GET", "/v1/threads", nil, params, "")
	if err != nil {
		return nil, err
	}
	return decodeJSON[ThreadList](data)
}

func (s *ThreadService) List(ctx context.Context, inboxID string, limit, offset int) (*ThreadList, error) {
	params := url.Values{"limit": {strconv.Itoa(limit)}, "offset": {strconv.Itoa(offset)}}
	data, err := s.c.do(ctx, "GET", fmt.Sprintf("/v1/inboxes/%s/threads", inboxID), nil, params, "")
	if err != nil {
		return nil, err
	}
	return decodeJSON[ThreadList](data)
}

func (s *ThreadService) Get(ctx context.Context, inboxID, threadID string) (*ThreadDetail, error) {
	data, err := s.c.do(ctx, "GET", fmt.Sprintf("/v1/inboxes/%s/threads/%s", inboxID, threadID), nil, nil, "")
	if err != nil {
		return nil, err
	}
	return decodeJSON[ThreadDetail](data)
}

func (s *ThreadService) GetAttachments(ctx context.Context, inboxID, threadID string) (*ThreadAttachments, error) {
	data, err := s.c.do(ctx, "GET", fmt.Sprintf("/v1/inboxes/%s/threads/%s/attachments", inboxID, threadID), nil, nil, "")
	if err != nil {
		return nil, err
	}
	return decodeJSON[ThreadAttachments](data)
}

// ── Webhooks ────────────────────────────────────────────────────────────────

type WebhookService struct{ c *Client }

type CreateWebhookParams struct {
	URL            string   `json:"url"`
	Events         []string `json:"events"`
	Description    string   `json:"description,omitempty"`
	IdempotencyKey string   `json:"-"`
}

func (s *WebhookService) Create(ctx context.Context, p *CreateWebhookParams) (*Webhook, error) {
	key := ""
	if p != nil {
		key = p.IdempotencyKey
	}
	data, err := s.c.do(ctx, "POST", "/v1/webhooks", p, nil, key)
	if err != nil {
		return nil, err
	}
	return decodeJSON[Webhook](data)
}

func (s *WebhookService) List(ctx context.Context) (*WebhookList, error) {
	data, err := s.c.do(ctx, "GET", "/v1/webhooks", nil, nil, "")
	if err != nil {
		return nil, err
	}
	return decodeJSON[WebhookList](data)
}

func (s *WebhookService) Delete(ctx context.Context, webhookID string) error {
	_, err := s.c.do(ctx, "DELETE", fmt.Sprintf("/v1/webhooks/%s", webhookID), nil, nil, "")
	return err
}

// ── Domains ─────────────────────────────────────────────────────────────────

type DomainService struct{ c *Client }

func (s *DomainService) Add(ctx context.Context, domain string) (*Domain, error) {
	data, err := s.c.do(ctx, "POST", "/v1/domains", map[string]string{"domain": domain}, nil, "")
	if err != nil {
		return nil, err
	}
	return decodeJSON[Domain](data)
}

func (s *DomainService) List(ctx context.Context) (*DomainList, error) {
	data, err := s.c.do(ctx, "GET", "/v1/domains", nil, nil, "")
	if err != nil {
		return nil, err
	}
	return decodeJSON[DomainList](data)
}

func (s *DomainService) Verify(ctx context.Context, domainID string) (*Domain, error) {
	data, err := s.c.do(ctx, "GET", fmt.Sprintf("/v1/domains/%s/verify", domainID), nil, nil, "")
	if err != nil {
		return nil, err
	}
	return decodeJSON[Domain](data)
}

func (s *DomainService) Delete(ctx context.Context, domainID string) error {
	_, err := s.c.do(ctx, "DELETE", fmt.Sprintf("/v1/domains/%s", domainID), nil, nil, "")
	return err
}

// ── API Keys ────────────────────────────────────────────────────────────────

type ApiKeyService struct{ c *Client }

func (s *ApiKeyService) Create(ctx context.Context, name string, permissions []string) (*ApiKey, error) {
	body := map[string]interface{}{"name": name}
	if len(permissions) > 0 {
		body["permissions"] = permissions
	}
	data, err := s.c.do(ctx, "POST", "/v1/api-keys", body, nil, "")
	if err != nil {
		return nil, err
	}
	return decodeJSON[ApiKey](data)
}

func (s *ApiKeyService) List(ctx context.Context) (*ApiKeyList, error) {
	data, err := s.c.do(ctx, "GET", "/v1/api-keys", nil, nil, "")
	if err != nil {
		return nil, err
	}
	return decodeJSON[ApiKeyList](data)
}

func (s *ApiKeyService) Delete(ctx context.Context, keyID string) error {
	_, err := s.c.do(ctx, "DELETE", fmt.Sprintf("/v1/api-keys/%s", keyID), nil, nil, "")
	return err
}

// ── Workspaces ──────────────────────────────────────────────────────────────

type WorkspaceService struct{ c *Client }

type CreateWorkspaceParams struct {
	Name           string `json:"name"`
	Description    string `json:"description,omitempty"`
	IdempotencyKey string `json:"-"`
}

type UpdateWorkspaceParams struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

func (s *WorkspaceService) Create(ctx context.Context, p *CreateWorkspaceParams) (*Workspace, error) {
	key := ""
	if p != nil {
		key = p.IdempotencyKey
	}
	data, err := s.c.do(ctx, "POST", "/v1/workspaces", p, nil, key)
	if err != nil {
		return nil, err
	}
	return decodeJSON[Workspace](data)
}

func (s *WorkspaceService) List(ctx context.Context) (*WorkspaceList, error) {
	data, err := s.c.do(ctx, "GET", "/v1/workspaces", nil, nil, "")
	if err != nil {
		return nil, err
	}
	return decodeJSON[WorkspaceList](data)
}

func (s *WorkspaceService) Get(ctx context.Context, workspaceID string) (*Workspace, error) {
	data, err := s.c.do(ctx, "GET", fmt.Sprintf("/v1/workspaces/%s", workspaceID), nil, nil, "")
	if err != nil {
		return nil, err
	}
	return decodeJSON[Workspace](data)
}

func (s *WorkspaceService) Update(ctx context.Context, workspaceID string, p *UpdateWorkspaceParams) (*Workspace, error) {
	data, err := s.c.do(ctx, "PUT", fmt.Sprintf("/v1/workspaces/%s", workspaceID), p, nil, "")
	if err != nil {
		return nil, err
	}
	return decodeJSON[Workspace](data)
}

func (s *WorkspaceService) Delete(ctx context.Context, workspaceID string) error {
	_, err := s.c.do(ctx, "DELETE", fmt.Sprintf("/v1/workspaces/%s", workspaceID), nil, nil, "")
	return err
}

// PodService is a deprecated alias for WorkspaceService.
type PodService = WorkspaceService

// ── Usage ───────────────────────────────────────────────────────────────────

type UsageService struct{ c *Client }

func (s *UsageService) Get(ctx context.Context) (*UsageInfo, error) {
	data, err := s.c.do(ctx, "GET", "/v1/usage", nil, nil, "")
	if err != nil {
		return nil, err
	}
	return decodeJSON[UsageInfo](data)
}

// ── Drafts ──────────────────────────────────────────────────────────────────

type DraftService struct{ c *Client }

type CreateDraftParams struct {
	To               []string `json:"to,omitempty"`
	CC               []string `json:"cc,omitempty"`
	BCC              []string `json:"bcc,omitempty"`
	Subject          string   `json:"subject,omitempty"`
	TextBody         string   `json:"text_body,omitempty"`
	HTMLBody         string   `json:"html_body,omitempty"`
	ReplyToMessageID string   `json:"reply_to_message_id,omitempty"`
}

func (s *DraftService) Create(ctx context.Context, inboxID string, p *CreateDraftParams) (*Draft, error) {
	data, err := s.c.do(ctx, "POST", fmt.Sprintf("/v1/inboxes/%s/drafts", inboxID), p, nil, "")
	if err != nil {
		return nil, err
	}
	return decodeJSON[Draft](data)
}

func (s *DraftService) List(ctx context.Context, inboxID string) (*DraftList, error) {
	data, err := s.c.do(ctx, "GET", fmt.Sprintf("/v1/inboxes/%s/drafts", inboxID), nil, nil, "")
	if err != nil {
		return nil, err
	}
	return decodeJSON[DraftList](data)
}

func (s *DraftService) Get(ctx context.Context, inboxID, draftID string) (*Draft, error) {
	data, err := s.c.do(ctx, "GET", fmt.Sprintf("/v1/inboxes/%s/drafts/%s", inboxID, draftID), nil, nil, "")
	if err != nil {
		return nil, err
	}
	return decodeJSON[Draft](data)
}

func (s *DraftService) Send(ctx context.Context, inboxID, draftID string) (*SendResult, error) {
	data, err := s.c.do(ctx, "POST", fmt.Sprintf("/v1/inboxes/%s/drafts/%s/send", inboxID, draftID), nil, nil, "")
	if err != nil {
		return nil, err
	}
	return decodeJSON[SendResult](data)
}

func (s *DraftService) Delete(ctx context.Context, inboxID, draftID string) error {
	_, err := s.c.do(ctx, "DELETE", fmt.Sprintf("/v1/inboxes/%s/drafts/%s", inboxID, draftID), nil, nil, "")
	return err
}

// ── Templates ───────────────────────────────────────────────────────────────

type TemplateService struct{ c *Client }

type CreateTemplateParams struct {
	Name     string `json:"name"`
	Subject  string `json:"subject,omitempty"`
	TextBody string `json:"text_body,omitempty"`
	HTMLBody string `json:"html_body,omitempty"`
}

func (s *TemplateService) Create(ctx context.Context, p *CreateTemplateParams) (*Template, error) {
	data, err := s.c.do(ctx, "POST", "/v1/templates", p, nil, "")
	if err != nil {
		return nil, err
	}
	return decodeJSON[Template](data)
}

func (s *TemplateService) List(ctx context.Context) (*TemplateList, error) {
	data, err := s.c.do(ctx, "GET", "/v1/templates", nil, nil, "")
	if err != nil {
		return nil, err
	}
	return decodeJSON[TemplateList](data)
}

func (s *TemplateService) Get(ctx context.Context, templateID string) (*Template, error) {
	data, err := s.c.do(ctx, "GET", fmt.Sprintf("/v1/templates/%s", templateID), nil, nil, "")
	if err != nil {
		return nil, err
	}
	return decodeJSON[Template](data)
}

func (s *TemplateService) Delete(ctx context.Context, templateID string) error {
	_, err := s.c.do(ctx, "DELETE", fmt.Sprintf("/v1/templates/%s", templateID), nil, nil, "")
	return err
}

// ── Filters ─────────────────────────────────────────────────────────────────

type FilterService struct{ c *Client }

type CreateFilterParams struct {
	Action      string `json:"action"`
	Criteria    string `json:"criteria"`
	Value       string `json:"value"`
	Description string `json:"description,omitempty"`
}

func (s *FilterService) Create(ctx context.Context, inboxID string, p *CreateFilterParams) (map[string]interface{}, error) {
	data, err := s.c.do(ctx, "POST", fmt.Sprintf("/v1/inboxes/%s/filters", inboxID), p, nil, "")
	if err != nil {
		return nil, err
	}
	out, err := decodeJSON[map[string]interface{}](data)
	if err != nil {
		return nil, err
	}
	return *out, nil
}

func (s *FilterService) List(ctx context.Context, inboxID string) (map[string]interface{}, error) {
	data, err := s.c.do(ctx, "GET", fmt.Sprintf("/v1/inboxes/%s/filters", inboxID), nil, nil, "")
	if err != nil {
		return nil, err
	}
	out, err := decodeJSON[map[string]interface{}](data)
	if err != nil {
		return nil, err
	}
	return *out, nil
}

func (s *FilterService) Delete(ctx context.Context, inboxID, filterID string) error {
	_, err := s.c.do(ctx, "DELETE", fmt.Sprintf("/v1/inboxes/%s/filters/%s", inboxID, filterID), nil, nil, "")
	return err
}

// ── Contacts ────────────────────────────────────────────────────────────────

type ContactService struct{ c *Client }

type CreateContactParams struct {
	Email    string                 `json:"email"`
	Name     string                 `json:"name,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

func (s *ContactService) Create(ctx context.Context, inboxID string, p *CreateContactParams) (*Contact, error) {
	data, err := s.c.do(ctx, "POST", fmt.Sprintf("/v1/inboxes/%s/contacts", inboxID), p, nil, "")
	if err != nil {
		return nil, err
	}
	return decodeJSON[Contact](data)
}

func (s *ContactService) List(ctx context.Context, inboxID string, q string) (*ContactList, error) {
	params := url.Values{}
	if q != "" {
		params.Set("q", q)
	}
	data, err := s.c.do(ctx, "GET", fmt.Sprintf("/v1/inboxes/%s/contacts", inboxID), nil, params, "")
	if err != nil {
		return nil, err
	}
	return decodeJSON[ContactList](data)
}

func (s *ContactService) Delete(ctx context.Context, inboxID, contactID string) error {
	_, err := s.c.do(ctx, "DELETE", fmt.Sprintf("/v1/inboxes/%s/contacts/%s", inboxID, contactID), nil, nil, "")
	return err
}

// ── Suppressions ────────────────────────────────────────────────────────────

type SuppressionService struct{ c *Client }

func (s *SuppressionService) Add(ctx context.Context, email, reason string) (map[string]interface{}, error) {
	body := map[string]string{"email": email, "reason": reason}
	data, err := s.c.do(ctx, "POST", "/v1/suppressions", body, nil, "")
	if err != nil {
		return nil, err
	}
	out, err := decodeJSON[map[string]interface{}](data)
	if err != nil {
		return nil, err
	}
	return *out, nil
}

func (s *SuppressionService) List(ctx context.Context) (map[string]interface{}, error) {
	data, err := s.c.do(ctx, "GET", "/v1/suppressions", nil, nil, "")
	if err != nil {
		return nil, err
	}
	out, err := decodeJSON[map[string]interface{}](data)
	if err != nil {
		return nil, err
	}
	return *out, nil
}

func (s *SuppressionService) Remove(ctx context.Context, email string) error {
	_, err := s.c.do(ctx, "DELETE", fmt.Sprintf("/v1/suppressions/%s", email), nil, nil, "")
	return err
}

// ── Events ──────────────────────────────────────────────────────────────────

type EventService struct{ c *Client }

func (s *EventService) List(ctx context.Context, eventType string, limit int, cursor string) (*EventList, error) {
	params := url.Values{}
	if eventType != "" {
		params.Set("type", eventType)
	}
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}
	if cursor != "" {
		params.Set("cursor", cursor)
	}
	data, err := s.c.do(ctx, "GET", "/v1/events", nil, params, "")
	if err != nil {
		return nil, err
	}
	return decodeJSON[EventList](data)
}

func (s *EventService) Get(ctx context.Context, eventID string) (*Event, error) {
	data, err := s.c.do(ctx, "GET", fmt.Sprintf("/v1/events/%s", eventID), nil, nil, "")
	if err != nil {
		return nil, err
	}
	return decodeJSON[Event](data)
}
