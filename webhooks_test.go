package deadsimple

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"
)

const (
	testSecret = "whsec_test_secret"
	testBody   = `{"event":"message.received","data":{"message_id":"msg_1"}}`
)

func sign(body string, ts int64, secret string) (string, string) {
	tsStr := fmt.Sprintf("%d", ts)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(tsStr + "." + body))
	return tsStr, hex.EncodeToString(mac.Sum(nil))
}

func testHeaders(body string, ts int64, both bool) http.Header {
	tsStr, digest := sign(body, ts, testSecret)
	h := http.Header{}
	h.Set("Content-Type", "application/json")
	h.Set("X-DSE-Webhook-Id", "wh_1")
	h.Set("X-DSE-Webhook-Timestamp", tsStr)
	h.Set("X-DSE-Webhook-Signature", "sha256="+digest)
	if both {
		h.Set("X-DSE-Signature", "t="+tsStr+",v1="+digest)
	}
	return h
}

func TestVerifyCombinedHeader(t *testing.T) {
	ts, digest := sign(testBody, time.Now().Unix(), testSecret)
	sig := "t=" + ts + ",v1=" + digest
	if err := VerifyWebhookSignature([]byte(testBody), sig, testSecret, 0); err != nil {
		t.Fatalf("expected valid signature, got %v", err)
	}
}

func TestVerifySplitHeaderWithTimestamp(t *testing.T) {
	ts, digest := sign(testBody, time.Now().Unix(), testSecret)
	if err := VerifyWebhookSignature([]byte(testBody), "sha256="+digest, testSecret, 0, ts); err != nil {
		t.Fatalf("expected valid signature, got %v", err)
	}
}

func TestSplitHeaderWithoutTimestampIsRejected(t *testing.T) {
	_, digest := sign(testBody, time.Now().Unix(), testSecret)
	err := VerifyWebhookSignature([]byte(testBody), "sha256="+digest, testSecret, 0)
	if err == nil || !strings.Contains(err.Error(), "X-DSE-Webhook-Timestamp") {
		t.Fatalf("expected a missing-timestamp error, got %v", err)
	}
}

func TestVerifyRequestReadsEitherHeader(t *testing.T) {
	now := time.Now().Unix()
	for _, both := range []bool{true, false} {
		if err := VerifyWebhookRequest([]byte(testBody), testHeaders(testBody, now, both), testSecret, 0); err != nil {
			t.Fatalf("both=%v: expected valid, got %v", both, err)
		}
	}
}

func TestTamperedBodyFails(t *testing.T) {
	h := testHeaders(testBody, time.Now().Unix(), true)
	if err := VerifyWebhookRequest([]byte(testBody+" "), h, testSecret, 0); err == nil {
		t.Fatal("expected a tampered body to fail verification")
	}
}

func TestWrongSecretFails(t *testing.T) {
	h := testHeaders(testBody, time.Now().Unix(), true)
	if err := VerifyWebhookRequest([]byte(testBody), h, "whsec_other", 0); err == nil {
		t.Fatal("expected the wrong secret to fail verification")
	}
}

func TestExpiredSignature(t *testing.T) {
	old := time.Now().Unix() - 3600
	h := testHeaders(testBody, old, true)

	err := VerifyWebhookRequest([]byte(testBody), h, testSecret, 0)
	if err == nil || !strings.Contains(err.Error(), "expired") {
		t.Fatalf("expected an expiry error, got %v", err)
	}
	if err := VerifyWebhookRequest([]byte(testBody), h, testSecret, 7200); err != nil {
		t.Fatalf("expected a wider tolerance to accept it, got %v", err)
	}
}

func TestMissingAndMalformedSignatures(t *testing.T) {
	if err := VerifyWebhookSignature([]byte(testBody), "", testSecret, 0); err == nil {
		t.Fatal("expected an empty signature to fail")
	}
	if err := VerifyWebhookSignature([]byte(testBody), "nonsense", testSecret, 0); err == nil {
		t.Fatal("expected a malformed signature to fail")
	}
	if err := VerifyWebhookSignature([]byte(testBody), "t=abc,v1=dead", testSecret, 0); err == nil {
		t.Fatal("expected a non-numeric timestamp to fail")
	}
	if err := VerifyWebhookRequest([]byte(testBody), http.Header{}, testSecret, 0); err == nil {
		t.Fatal("expected a request with no signature header to fail")
	}
}
