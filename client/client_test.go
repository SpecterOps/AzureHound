package client

import (
	"context"
	"io"
	"net/http"
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
		require.True(t, ok, "expected a value on out, channel closed")
		require.Error(t, result.Error, "expected an error result from timed-out request")
		require.ErrorIs(t, result.Error, context.DeadlineExceeded)
	case <-time.After(5 * time.Second):
		t.Fatal("getAzureObjectList did not return within expected timeout; pipeline is hung")
	}

	// drain and ensure channel closes
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
