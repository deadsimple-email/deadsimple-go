package deadsimple

// Inbox represents an email inbox.
type Inbox struct {
	InboxID      string   `json:"inbox_id"`
	Email        string   `json:"email"`
	LocalPart    string   `json:"local_part"`
	Domain       string   `json:"domain"`
	DisplayName  string   `json:"display_name"`
	Tags         []string `json:"tags"`
	Status       string   `json:"status"`
	MessageCount int      `json:"message_count"`
	CreatedAt    string   `json:"created_at"`
	UpdatedAt    string   `json:"updated_at"`
}

// InboxList is a paginated list of inboxes.
type InboxList struct {
	Inboxes []Inbox `json:"inboxes"`
	Count   int     `json:"count"`
	Cursor  string  `json:"cursor,omitempty"`
}

// BulkCreateResult contains results of a bulk inbox creation.
type BulkCreateResult struct {
	Results []map[string]interface{} `json:"results"`
	Created int                      `json:"created"`
	Failed  int                      `json:"failed"`
	Total   int                      `json:"total"`
}

// Message represents an email message.
type Message struct {
	MessageID   string       `json:"message_id"`
	InboxID     string       `json:"inbox_id"`
	ThreadID    string       `json:"thread_id"`
	Direction   string       `json:"direction"`
	FromEmail   string       `json:"from_email"`
	FromName    string       `json:"from_name"`
	To          []string     `json:"to"`
	CC          []string     `json:"cc"`
	Subject     string       `json:"subject"`
	Labels      []string     `json:"labels"`
	Status      string       `json:"status"`
	IsSpam      bool         `json:"is_spam,omitempty"`
	SpamScore   float64      `json:"spam_score,omitempty"`
	Snippet     string       `json:"snippet,omitempty"`
	Attachments []Attachment `json:"attachments"`
	CreatedAt   string       `json:"created_at"`
	ReceivedAt  string       `json:"received_at"`
	// Full body fields (only when fetching single message)
	TextBody    string            `json:"text_body,omitempty"`
	HTMLBody    string            `json:"html_body,omitempty"`
	LatestReply string            `json:"latest_reply,omitempty"`
	InReplyTo   string            `json:"in_reply_to,omitempty"`
	References  []string          `json:"references,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
	// Extracted data
	Extracted map[string]interface{} `json:"extracted,omitempty"`
	// Injection risk
	InjectionRisk  string  `json:"injection_risk,omitempty"`
	InjectionScore float64 `json:"injection_score,omitempty"`
}

// Attachment describes a file attached to a message.
type Attachment struct {
	AttachmentID string `json:"attachment_id"`
	Filename     string `json:"filename"`
	ContentType  string `json:"content_type"`
	Size         int    `json:"size"`
}

// MessageList is a paginated list of messages.
type MessageList struct {
	Messages []Message `json:"messages"`
	Count    int       `json:"count"`
	Cursor   string    `json:"cursor,omitempty"`
}

// SendResult is returned after sending a message.
type SendResult struct {
	MessageID string `json:"message_id"`
	Status    string `json:"status"`
	SendAt    string `json:"send_at,omitempty"`
}

// Thread represents a conversation thread.
type Thread struct {
	ThreadID      string                 `json:"thread_id"`
	Subject       string                 `json:"subject"`
	MessageCount  int                    `json:"message_count"`
	Participants  []string               `json:"participants"`
	LatestMessage map[string]interface{} `json:"latest_message"`
	CreatedAt     string                 `json:"created_at"`
	UpdatedAt     string                 `json:"updated_at"`
}

// ThreadList is a paginated list of threads.
type ThreadList struct {
	Threads []Thread `json:"threads"`
	Count   int      `json:"count"`
	Total   int      `json:"total"`
}

// ThreadDetail contains all messages in a thread.
type ThreadDetail struct {
	ThreadID     string    `json:"thread_id"`
	Subject      string    `json:"subject"`
	Messages     []Message `json:"messages"`
	Participants []string  `json:"participants"`
}

// ThreadAttachmentItem represents a single attachment with message attribution.
type ThreadAttachmentItem struct {
	AttachmentID string `json:"attachment_id"`
	MessageID    string `json:"message_id"`
	Filename     string `json:"filename"`
	ContentType  string `json:"content_type"`
	Size         int    `json:"size"`
}

// ThreadAttachments contains all attachments across a thread.
type ThreadAttachments struct {
	ThreadID    string                 `json:"thread_id"`
	Attachments []ThreadAttachmentItem `json:"attachments"`
	Count       int                    `json:"count"`
}

// Webhook event type constants.
const (
	EventMessageReceived    = "message.received"
	EventMessageSent        = "message.sent"
	EventMessageDelivered   = "message.delivered"
	EventMessageBounced     = "message.bounced"
	EventMessageComplained  = "message.complained"
	EventMessageRejected    = "message.rejected"
	EventMessageOpened      = "message.opened"
	EventMessageLinkClicked = "message.link_clicked"
	EventInboxCreated       = "inbox.created"
	EventInboxDeleted       = "inbox.deleted"
	EventDomainVerified     = "domain.verified"
)

// AllWebhookEvents is the complete list of available webhook event types.
var AllWebhookEvents = []string{
	EventMessageReceived, EventMessageSent, EventMessageDelivered,
	EventMessageBounced, EventMessageComplained, EventMessageRejected,
	EventMessageOpened, EventMessageLinkClicked,
	EventInboxCreated, EventInboxDeleted, EventDomainVerified,
}

// Webhook represents a webhook endpoint.
type Webhook struct {
	WebhookID       string   `json:"webhook_id"`
	URL             string   `json:"url"`
	Events          []string `json:"events"`
	Description     string   `json:"description"`
	IsActive        bool     `json:"is_active"`
	CreatedAt       string   `json:"created_at"`
	LastTriggeredAt string   `json:"last_triggered_at,omitempty"`
	FailureCount    int      `json:"failure_count"`
	SigningSecret   string   `json:"signing_secret,omitempty"`
	// Headers are the custom delivery headers. Values are masked by the API
	// (write-only).
	Headers map[string]string `json:"headers,omitempty"`
}

// WebhookList is a list of webhooks.
type WebhookList struct {
	Webhooks []Webhook `json:"webhooks"`
	Count    int       `json:"count"`
}

// Domain represents a custom domain.
type Domain struct {
	DomainID   string              `json:"domain_id"`
	DomainName string              `json:"domain"`
	Status     string              `json:"status"`
	DNSRecords []map[string]string `json:"dns_records"`
	CreatedAt  string              `json:"created_at"`
	VerifiedAt string              `json:"verified_at,omitempty"`
}

// DomainList is a list of domains.
type DomainList struct {
	Domains []Domain `json:"domains"`
	Count   int      `json:"count"`
}

// ApiKey represents an API key.
type ApiKey struct {
	KeyID       string   `json:"key_id"`
	Name        string   `json:"name"`
	KeyPrefix   string   `json:"key_prefix"`
	Permissions []string `json:"permissions"`
	IsActive    bool     `json:"is_active"`
	CreatedAt   string   `json:"created_at"`
	Key         string   `json:"key,omitempty"`
	RevokedAt   string   `json:"revoked_at,omitempty"`
}

// ApiKeyList is a list of API keys.
type ApiKeyList struct {
	ApiKeys []ApiKey `json:"api_keys"`
	Count   int      `json:"count"`
}

// Workspace represents a multi-tenant workspace (formerly pod).
type Workspace struct {
	WorkspaceID string            `json:"workspace_id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	IsDefault   bool              `json:"is_default"`
	CreatedAt   string            `json:"created_at"`
	UpdatedAt   string            `json:"updated_at"`
	InboxCount  int               `json:"inbox_count"`
	ApiKey      map[string]string `json:"api_key,omitempty"`
}

// Pod is a deprecated alias for Workspace.
type Pod = Workspace

// WorkspaceList is a list of workspaces.
type WorkspaceList struct {
	Workspaces []Workspace `json:"workspaces"`
	Count      int         `json:"count"`
}

// PodList is a deprecated alias for WorkspaceList.
type PodList = WorkspaceList

// UsageInfo contains account usage metrics.
type UsageInfo struct {
	Plan          string         `json:"plan"`
	PlanName      string         `json:"plan_name"`
	Inboxes       map[string]int `json:"inboxes"`
	Emails        map[string]int `json:"emails"`
	Webhooks      int            `json:"webhooks"`
	CustomDomains map[string]int `json:"custom_domains"`
}

// Draft represents a draft email.
type Draft struct {
	DraftID          string   `json:"draft_id"`
	InboxID          string   `json:"inbox_id"`
	To               []string `json:"to"`
	CC               []string `json:"cc"`
	BCC              []string `json:"bcc"`
	Subject          string   `json:"subject"`
	TextBody         string   `json:"text_body"`
	HTMLBody         string   `json:"html_body"`
	ReplyToMessageID string   `json:"reply_to_message_id,omitempty"`
	CreatedAt        string   `json:"created_at"`
	UpdatedAt        string   `json:"updated_at"`
}

// DraftList is a list of drafts.
type DraftList struct {
	Drafts []Draft `json:"drafts"`
	Count  int     `json:"count"`
}

// Template represents an email template.
type Template struct {
	TemplateID string   `json:"template_id"`
	Name       string   `json:"name"`
	Subject    string   `json:"subject"`
	TextBody   string   `json:"text_body"`
	HTMLBody   string   `json:"html_body"`
	Variables  []string `json:"variables"`
	Version    int      `json:"version"`
	CreatedAt  string   `json:"created_at"`
	UpdatedAt  string   `json:"updated_at"`
}

// TemplateList is a list of templates.
type TemplateList struct {
	Templates []Template `json:"templates"`
	Count     int        `json:"count"`
}

// Event represents an audit log entry.
type Event struct {
	EventID     string                 `json:"event_id"`
	EventType   string                 `json:"event_type"`
	Description string                 `json:"description"`
	Payload     map[string]interface{} `json:"payload"`
	CreatedAt   string                 `json:"created_at"`
}

// EventList is a paginated list of events.
type EventList struct {
	Events []Event `json:"events"`
	Count  int     `json:"count"`
	Cursor string  `json:"cursor,omitempty"`
}

// Contact represents a contact in the address book.
type Contact struct {
	ContactID       string                 `json:"contact_id"`
	InboxID         string                 `json:"inbox_id"`
	Email           string                 `json:"email"`
	Name            string                 `json:"name"`
	Metadata        map[string]interface{} `json:"metadata"`
	MessageCount    int                    `json:"message_count"`
	LastContactedAt string                 `json:"last_contacted_at,omitempty"`
	CreatedAt       string                 `json:"created_at"`
	UpdatedAt       string                 `json:"updated_at"`
}

// ContactList is a list of contacts.
type ContactList struct {
	Contacts []Contact `json:"contacts"`
	Count    int       `json:"count"`
}
