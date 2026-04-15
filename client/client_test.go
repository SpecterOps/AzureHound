package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/bloodhoundad/azurehound/v2/client/query"
	"github.com/stretchr/testify/require"
)

// fakeRestClient is a minimal test double for rest.RestClient that allows
// controlling the Get response per test case.
type fakeRestClient struct {
	getFunc func(ctx context.Context, path string, params query.Params, headers map[string]string) (*http.Response, error)
}

func (s *fakeRestClient) Get(ctx context.Context, path string, params query.Params, headers map[string]string) (*http.Response, error) {
	return s.getFunc(ctx, path, params, headers)
}

func (s *fakeRestClient) Delete(context.Context, string, interface{}, query.Params, map[string]string) (*http.Response, error) {
	return nil, nil
}
func (s *fakeRestClient) Patch(context.Context, string, interface{}, query.Params, map[string]string) (*http.Response, error) {
	return nil, nil
}
func (s *fakeRestClient) Post(context.Context, string, interface{}, query.Params, map[string]string) (*http.Response, error) {
	return nil, nil
}
func (s *fakeRestClient) Put(context.Context, string, interface{}, query.Params, map[string]string) (*http.Response, error) {
	return nil, nil
}
func (s *fakeRestClient) Send(req *http.Request) (*http.Response, error) { return nil, nil }
func (s *fakeRestClient) AddAuthenticationToRequest(req *http.Request) (*http.Request, error) {
	return req, nil
}
func (s *fakeRestClient) CloseIdleConnections() {}

func TestGetAzureObjectList_SuccessfulResponse(t *testing.T) {
	body := `{"value": [{"id": "1"}, {"id": "2"}]}`
	client := &fakeRestClient{
		getFunc: func(ctx context.Context, path string, params query.Params, headers map[string]string) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(body)),
			}, nil
		},
	}

	out := make(chan AzureResult[map[string]string])
	go getAzureObjectList(client, context.Background(), "/test/path", nil, out)

	var results []map[string]string
	for result := range out {
		require.NoError(t, result.Error)
		results = append(results, result.Ok)
	}

	require.Len(t, results, 2)
	require.Equal(t, "1", results[0]["id"])
	require.Equal(t, "2", results[1]["id"])
}

func TestGetAzureObjectList_HungResponseTimesOut(t *testing.T) {
	// Shorten the timeout so the test completes quickly
	original := pageRequestTimeout
	pageRequestTimeout = 500 * time.Millisecond
	defer func() { pageRequestTimeout = original }()

	client := &fakeRestClient{
		getFunc: func(ctx context.Context, path string, params query.Params, headers map[string]string) (*http.Response, error) {
			// Verify the context has a deadline (set by pageRequestTimeout)
			_, hasDeadline := ctx.Deadline()
			require.True(t, hasDeadline, "expected context passed to Get to have a deadline from pageRequestTimeout")

			// Simulate a hung connection: block until the context expires
			<-ctx.Done()
			return nil, ctx.Err()
		},
	}

	out := make(chan AzureResult[map[string]string])
	go getAzureObjectList(client, context.Background(), "/test/path", nil, out)

	// The channel should produce an error and close well within a few seconds
	select {
	case result, ok := <-out:
		if ok {
			require.Error(t, result.Error, "expected an error result from timed-out request")
			require.ErrorIs(t, result.Error, context.DeadlineExceeded)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("getAzureObjectList did not return within expected timeout; pipeline is hung")
	}

	// Drain and ensure channel closes
	for range out {
	}
}

func TestGetAzureObjectList_ParentContextCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	client := &fakeRestClient{
		getFunc: func(ctx context.Context, path string, params query.Params, headers map[string]string) (*http.Response, error) {
			<-ctx.Done()
			return nil, ctx.Err()
		},
	}

	out := make(chan AzureResult[map[string]string])
	go getAzureObjectList(client, ctx, "/test/path", nil, out)

	// Cancel the parent context after a short delay
	time.AfterFunc(100*time.Millisecond, cancel)

	select {
	case result, ok := <-out:
		if ok {
			require.Error(t, result.Error)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("getAzureObjectList did not respect parent context cancellation")
	}

	for range out {
	}
}

// TestGetAzureObjectList_StalledResponseBody_Integration starts a real HTTP
// server that sends response headers and a partial JSON body, then stalls
// indefinitely. This reproduces the exact failure mode from BED-4600: the
// server responds (so ResponseHeaderTimeout doesn't help) but the body read
// hangs, blocking the collection pipeline. The test verifies that the
// per-page context timeout terminates the hung read.
func TestGetAzureObjectList_StalledResponseBody_Integration(t *testing.T) {
	original := pageRequestTimeout
	pageRequestTimeout = 500 * time.Millisecond
	defer func() { pageRequestTimeout = original }()

	serverStalled := make(chan struct{})

	// Start a test server that sends headers + partial body, then blocks
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Send the beginning of a valid JSON response, but don't finish it
		fmt.Fprint(w, `{"value": [{"id": "1"}`)
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}

		// Block until the test completes (simulates a stalled connection)
		<-serverStalled
	}))
	defer server.Close()
	defer close(serverStalled)

	// Create a fakeRestClient that makes a real HTTP call to the stalling server.
	// This exercises the actual HTTP transport, TCP connection, and body read path.
	httpClient := server.Client()
	client := &fakeRestClient{
		getFunc: func(ctx context.Context, path string, params query.Params, headers map[string]string) (*http.Response, error) {
			req, err := http.NewRequestWithContext(ctx, "GET", server.URL+path, nil)
			if err != nil {
				return nil, err
			}
			return httpClient.Do(req)
		},
	}

	out := make(chan AzureResult[map[string]string])
	go getAzureObjectList(client, context.Background(), "/test/path", nil, out)

	select {
	case result, ok := <-out:
		if ok {
			require.Error(t, result.Error, "expected an error from stalled response body")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("getAzureObjectList hung on stalled response body; timeout did not fire")
	}

	for range out {
	}
}
