package gobalena

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

const testDeviceUUID = "2305778e615ed9ac3652699962bf3cf9"

// newCapturingServer returns a test server that records the raw request body of the last request it
// received and always answers 200, plus a pointer to that captured body.
func newCapturingServer(t *testing.T) (*httptest.Server, *[]byte) {
	t.Helper()

	var captured []byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("reading request body: %v", err)
		}
		captured = body
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)

	return server, &captured
}

// TestForceApplySendsForceUnderDataKey pins the supervisor-proxy contract: force must travel under
// "data". The proxy silently drops unknown keys and still answers 200, so sending it anywhere else
// downgrades the call to a non-forced update that honours update lockfiles - with no error to show
// for it. This asserts the serialized payload rather than the map handed to resty, because only the
// serialized form shows what the proxy actually reads.
func TestForceApplySendsForceUnderDataKey(t *testing.T) {
	server, captured := newCapturingServer(t)

	if err := NewCloudClient("test-key", server.URL).ForceApply(context.Background(), testDeviceUUID); err != nil {
		t.Fatalf("ForceApply returned an unexpected error: %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(*captured, &payload); err != nil {
		t.Fatalf("unmarshalling captured request body %q: %v", string(*captured), err)
	}

	data, ok := payload["data"].(map[string]any)
	if !ok {
		t.Fatalf(`expected a "data" object in the payload, got: %s`, string(*captured))
	}

	if force, ok := data["force"].(bool); !ok || !force {
		t.Errorf(`expected data.force to be true, got: %s`, string(*captured))
	}

	if _, exists := payload["body"]; exists {
		t.Errorf(`payload must not carry a "body" key; the proxy ignores it and force is lost: %s`, string(*captured))
	}

	if uuid, _ := payload["uuid"].(string); uuid != testDeviceUUID {
		t.Errorf("expected uuid %q, got %q", testDeviceUUID, uuid)
	}
}
