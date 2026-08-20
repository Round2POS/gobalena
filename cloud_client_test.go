package gobalena

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"testing/iotest"
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

// newDownloadOSServer answers the fleet lookup DownloadOS makes first, then hands /download to the
// supplied handler, so a test only has to describe the download response it cares about.
func newDownloadOSServer(t *testing.T, download http.HandlerFunc) *httptest.Server {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/download" {
			download(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte(`{"d":[{"id":2076891,"app_name":"activation"}]}`)); err != nil {
			t.Errorf("writing fleet response: %v", err)
		}
	}))
	t.Cleanup(server.Close)

	return server
}

// TestDownloadOSErrorNamesStatusAndBody pins what a failed download reports. DownloadOS sets
// SetDoNotParseResponse so it can stream a multi-GB image, which leaves response.Body() empty - so
// building the error from it produced "error downloading os:" and nothing more, and dropped the
// status. A 404 from a bad OS version then read exactly like an expired token, which is how a
// broken production download path stayed invisible. Each case asserts both halves are present.
func TestDownloadOSErrorNamesStatusAndBody(t *testing.T) {
	tests := []struct {
		name   string
		status int
		body   string
	}{
		{name: "not found", status: http.StatusNotFound, body: `{"error":"no such version"}`},
		{name: "unauthorized", status: http.StatusUnauthorized, body: `{"error":"token expired"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := newDownloadOSServer(t, func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tt.status)
				if _, err := w.Write([]byte(tt.body)); err != nil {
					t.Errorf("writing download response: %v", err)
				}
			})

			_, err := NewCloudClient("test-key", server.URL).DownloadOS(
				context.Background(), io.Discard, "activation", DeviceTypeAMD64, "6.12.2", nil,
			)
			if err == nil {
				t.Fatal("expected an error from a failed download, got nil")
			}

			if !strings.Contains(err.Error(), strconv.Itoa(tt.status)) {
				t.Errorf("error must name the HTTP status %d, got: %v", tt.status, err)
			}

			if !strings.Contains(err.Error(), tt.body) {
				t.Errorf("error must carry balena's response body %q, got: %v", tt.body, err)
			}
		})
	}
}

// TestDownloadOSErrorsDifferByStatus pins that two failures are distinguishable. The bug was not
// only a missing status, it was that every failure produced the identical string.
func TestDownloadOSErrorsDifferByStatus(t *testing.T) {
	message := func(t *testing.T, status int) string {
		t.Helper()

		server := newDownloadOSServer(t, func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(status)
		})

		_, err := NewCloudClient("test-key", server.URL).DownloadOS(
			context.Background(), io.Discard, "activation", DeviceTypeAMD64, "6.12.2", nil,
		)
		if err == nil {
			t.Fatalf("expected an error for status %d, got nil", status)
		}

		return err.Error()
	}

	notFound := message(t, http.StatusNotFound)
	if unauthorized := message(t, http.StatusUnauthorized); notFound == unauthorized {
		t.Errorf("a 404 and a 401 must not read identically, both were: %s", notFound)
	}
}

// TestDownloadOSStreamsSuccessfulResponse pins that reading the error body did not disturb the
// success path, which must still stream straight to the writer and answer the balena filename.
func TestDownloadOSStreamsSuccessfulResponse(t *testing.T) {
	const image = "not-really-an-os-image"

	server := newDownloadOSServer(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Disposition", `attachment; filename=balena-cloud-activation.zip`)
		if _, err := w.Write([]byte(image)); err != nil {
			t.Errorf("writing download response: %v", err)
		}
	})

	var downloaded bytes.Buffer
	filename, err := NewCloudClient("test-key", server.URL).DownloadOS(
		context.Background(), &downloaded, "activation", DeviceTypeAMD64, "6.12.2", nil,
	)
	if err != nil {
		t.Fatalf("DownloadOS returned an unexpected error: %v", err)
	}

	if downloaded.String() != image {
		t.Errorf("expected the image streamed through verbatim, got: %q", downloaded.String())
	}

	if filename != "activation.zip" {
		t.Errorf(`expected filename "activation.zip", got: %q`, filename)
	}
}

// TestDownloadOSCapsAndMarksAnOversizeErrorBody pins the bound on the error body read. resty's own
// SetResponseBodyLimit is not enforced when DoNotParseResponse is set, so nothing but this cap
// stops a pathological response being read whole into a message. Cutting it silently would be its
// own trap - a truncated JSON body reads exactly like a complete one - so the cut is marked.
func TestDownloadOSCapsAndMarksAnOversizeErrorBody(t *testing.T) {
	const tail = "TAIL-THAT-MUST-NOT-APPEAR"
	oversize := strings.Repeat("A", maxErrorBodyBytes*2) + tail

	server := newDownloadOSServer(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		if _, err := w.Write([]byte(oversize)); err != nil {
			t.Errorf("writing download response: %v", err)
		}
	})

	_, err := NewCloudClient("test-key", server.URL).DownloadOS(
		context.Background(), io.Discard, "activation", DeviceTypeAMD64, "6.12.2", nil,
	)
	if err == nil {
		t.Fatal("expected an error from a failed download, got nil")
	}

	if strings.Contains(err.Error(), tail) {
		t.Error("the body was read past the cap; the whole response reached the message")
	}

	if len(err.Error()) > maxErrorBodyBytes*2 {
		t.Errorf("message is unbounded at %d bytes for a %d byte body", len(err.Error()), len(oversize))
	}

	if !strings.HasSuffix(err.Error(), errorBodyTruncated) {
		t.Errorf("a cut body must say so, got the tail: %q", err.Error()[len(err.Error())-40:])
	}
}

// TestDownloadOSDoesNotMarkAWholeErrorBody is the other half of the cap: a body that fits must not
// be labelled truncated, or the marker means nothing and every message carries it.
func TestDownloadOSDoesNotMarkAWholeErrorBody(t *testing.T) {
	server := newDownloadOSServer(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		if _, err := w.Write([]byte(strings.Repeat("A", maxErrorBodyBytes))); err != nil {
			t.Errorf("writing download response: %v", err)
		}
	})

	_, err := NewCloudClient("test-key", server.URL).DownloadOS(
		context.Background(), io.Discard, "activation", DeviceTypeAMD64, "6.12.2", nil,
	)
	if err == nil {
		t.Fatal("expected an error from a failed download, got nil")
	}

	if strings.Contains(err.Error(), errorBodyTruncated) {
		t.Error("a body that fits the cap exactly must not be marked truncated")
	}
}

// countingReader reports how many bytes were actually pulled from the source, which is the only way
// to see the read bound. Every message-level assertion is blind to it: the body is capped for
// display either way, so a missing io.LimitReader would buffer the whole response unnoticed.
type countingReader struct {
	source io.Reader
	read   int
}

func (c *countingReader) Read(p []byte) (int, error) {
	n, err := c.source.Read(p)
	c.read += n

	return n, err
}

// TestReadErrorBodyBoundsBothTheReadAndTheResult pins the two halves of the cap separately. The
// returned body must not exceed maxErrorBodyBytes, so the error message honours the documented
// limit. The read must not exceed one byte past it, so a pathological response is never buffered
// whole just to build a message. Only the lookahead byte can tell a cut body from a whole one, so
// the two bounds differ by exactly one and neither implies the other.
func TestReadErrorBodyBoundsBothTheReadAndTheResult(t *testing.T) {
	source := &countingReader{source: strings.NewReader(strings.Repeat("A", maxErrorBodyBytes*25))}

	body, cut, err := readErrorBody(source)
	if err != nil {
		t.Fatalf("readErrorBody returned an unexpected error: %v", err)
	}

	if !cut {
		t.Error("a body far over the cap must report as cut")
	}

	if len(body) != maxErrorBodyBytes {
		t.Errorf("returned body must be capped at %d bytes, got %d", maxErrorBodyBytes, len(body))
	}

	if source.read > maxErrorBodyBytes+1 {
		t.Errorf("read %d bytes; must stop one byte past the cap, at %d", source.read, maxErrorBodyBytes+1)
	}
}

// TestReadErrorBodyKeepsABodyThatFits is the boundary case. A body of exactly the cap is whole, not
// cut, and must come back untouched - otherwise the truncation marker fires on complete responses.
func TestReadErrorBodyKeepsABodyThatFits(t *testing.T) {
	whole := strings.Repeat("A", maxErrorBodyBytes)

	body, cut, err := readErrorBody(strings.NewReader(whole))
	if err != nil {
		t.Fatalf("readErrorBody returned an unexpected error: %v", err)
	}

	if cut {
		t.Error("a body of exactly the cap must not report as cut")
	}

	if string(body) != whole {
		t.Errorf("expected the body returned whole, got %d of %d bytes", len(body), len(whole))
	}
}

// TestReadErrorBodyReportsAFailedRead covers the path where the body dies mid-read, which must
// surface rather than pass an empty body off as the response.
func TestReadErrorBodyReportsAFailedRead(t *testing.T) {
	_, _, err := readErrorBody(iotest.ErrReader(errors.New("connection reset")))

	if err == nil {
		t.Fatal("expected the read error to surface, got nil")
	}
}
