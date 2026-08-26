package deadsimple

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Signature headers sent with every webhook delivery.
const (
	signatureHeader        = "X-DSE-Signature"
	webhookSignatureHeader = "X-DSE-Webhook-Signature"
	webhookTimestampHeader = "X-DSE-Webhook-Timestamp"
)

// VerifyWebhookSignature verifies an HMAC-SHA256 webhook signature.
//
// payload is the raw request body, before any JSON decoding. signature is either
// the X-DSE-Signature header ("t=<unix>,v1=<hex>") or the X-DSE-Webhook-Signature
// header ("sha256=<hex>"); the latter carries no timestamp, so pass the
// X-DSE-Webhook-Timestamp value as timestamp. secret is the per-webhook signing
// secret. tolerance is the maximum age in seconds (use 0 for the default of 300).
//
// Prefer VerifyWebhookRequest, which reads whichever headers are present.
func VerifyWebhookSignature(payload []byte, signature, secret string, tolerance int, timestamp ...string) error {
	if signature == "" {
		return fmt.Errorf("deadsimple: missing signature header")
	}
	if tolerance <= 0 {
		tolerance = 300
	}

	var suppliedTimestamp string
	if len(timestamp) > 0 {
		suppliedTimestamp = strings.TrimSpace(timestamp[0])
	}

	tsStr, expectedSig, err := parseWebhookSignature(strings.TrimSpace(signature), suppliedTimestamp)
	if err != nil {
		return err
	}

	ts, err := strconv.ParseInt(tsStr, 10, 64)
	if err != nil {
		return fmt.Errorf("deadsimple: invalid timestamp in signature")
	}

	// Check freshness
	age := math.Abs(float64(time.Now().Unix() - ts))
	if age > float64(tolerance) {
		return fmt.Errorf("deadsimple: signature expired (%ds old, tolerance is %ds)", int(age), tolerance)
	}

	// Compute HMAC-SHA256 of "timestamp.payload"
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(fmt.Sprintf("%d.", ts)))
	mac.Write(payload)
	computed := hex.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(computed), []byte(expectedSig)) {
		return fmt.Errorf("deadsimple: signature mismatch")
	}

	return nil
}

// VerifyWebhookRequest verifies a webhook using a request's raw body and headers.
//
// It reads whichever signature header is present, so it works regardless of which
// form the delivery used. Read the body yourself first: http.Request.Body is a
// stream, and verification needs the exact bytes that were signed.
//
//	body, _ := io.ReadAll(r.Body)
//	if err := deadsimple.VerifyWebhookRequest(body, r.Header, secret, 0); err != nil {
//	    http.Error(w, "bad signature", http.StatusUnauthorized)
//	    return
//	}
func VerifyWebhookRequest(payload []byte, headers http.Header, secret string, tolerance int) error {
	if combined := headers.Get(signatureHeader); combined != "" {
		return VerifyWebhookSignature(payload, combined, secret, tolerance)
	}
	if split := headers.Get(webhookSignatureHeader); split != "" {
		return VerifyWebhookSignature(payload, split, secret, tolerance, headers.Get(webhookTimestampHeader))
	}
	return fmt.Errorf("deadsimple: missing signature header")
}

// parseWebhookSignature returns (timestamp, hex digest) from either encoding.
func parseWebhookSignature(signature, suppliedTimestamp string) (string, string, error) {
	// "sha256=<hex>", paired with a separate timestamp header.
	if strings.HasPrefix(signature, "sha256=") {
		if suppliedTimestamp == "" {
			return "", "", fmt.Errorf("deadsimple: a sha256= signature needs the %s header; pass it as timestamp or use VerifyWebhookRequest", webhookTimestampHeader)
		}
		return suppliedTimestamp, strings.TrimSpace(strings.TrimPrefix(signature, "sha256=")), nil
	}

	// "t=<unix>,v1=<hex>"
	parts := map[string]string{}
	for _, segment := range strings.Split(signature, ",") {
		if idx := strings.Index(segment, "="); idx > 0 {
			parts[strings.TrimSpace(segment[:idx])] = strings.TrimSpace(segment[idx+1:])
		}
	}

	ts, ok := parts["t"]
	if !ok || ts == "" {
		ts = suppliedTimestamp
	}
	expectedSig := parts["v1"]
	if ts == "" || expectedSig == "" {
		return "", "", fmt.Errorf("deadsimple: malformed signature header")
	}
	return ts, expectedSig, nil
}
