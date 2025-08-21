package bloodhound

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/http2"
)

func TestBHEClient_SendRequest(t *testing.T) {
	t.Run("GOAWAY error handling", func(t *testing.T) {
		client := &BHEClient{
			httpClient: &http.Client{
				Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
					return nil, &http2.GoAwayError{
						LastStreamID: 1,
						ErrCode:      http2.ErrCodeNo,
						DebugData:    "",
					}
				}),
			},
			maxRetries:   0,
			log:          logr.Discard(),
			requestLimit: 10,
		}

		req, _ := http.NewRequest("GET", "http://bhe.com", nil)
		_, err := client.SendRequest(req)

		require.Error(t, err)
	})

	t.Run("retry after failures", func(t *testing.T) {
		requestCount := 0
		maxRetries := 5

		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestCount++
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer testServer.Close()

		testUrl, _ := url.Parse(testServer.URL)

		client := &BHEClient{
			httpClient: http.DefaultClient,
			maxRetries: maxRetries,
			retryDelay: 0,
		}

		req, _ := http.NewRequest("GET", testUrl.String(), nil)
		_, err := client.SendRequest(req)

		require.Error(t, err)
		require.Equal(t, maxRetries+1, requestCount)
	})
}

func TestBHEClient_Ingest(t *testing.T) {
	t.Run("successful ingest request", func(t *testing.T) {
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusAccepted)
		}))
		defer testServer.Close()

		testUrl, _ := url.Parse(testServer.URL)

		client, err := NewBHEClient(*testUrl, "tokenId", "token", "", 1, 1, logr.Logger{})
		require.NoError(t, err)

		data := make(chan []any, 1)
		data <- []any{"test"}
		close(data)

		hadErrors := client.Ingest(context.Background(), data)

		require.False(t, hadErrors)
	})

	t.Run("retry after failures", func(t *testing.T) {
		requestCount := 0
		maxRetries := 1
		wg := sync.WaitGroup{}
		wg.Add(maxRetries + 1)

		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer wg.Done()
			requestCount++
			w.WriteHeader(http.StatusGatewayTimeout)
		}))
		defer testServer.Close()

		testUrl, _ := url.Parse(testServer.URL)

		client := &BHEClient{
			httpClient: http.DefaultClient,
			maxRetries: maxRetries,
			retryDelay: 0,
			bheUrl:     *testUrl,
		}
		data := make(chan []any, 1)
		data <- []any{"test"}
		close(data)

		hadErrors := client.Ingest(context.Background(), data)

		wg.Wait()
		require.True(t, hadErrors)
		require.Equal(t, maxRetries+1, requestCount)
	})

	t.Run("GOAWAY error handling", func(t *testing.T) {
		client := &BHEClient{
			httpClient: &http.Client{
				Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
					return nil, &http2.GoAwayError{
						ErrCode:   http2.ErrCodeNo,
						DebugData: "",
					}
				}),
			},
			log: logr.Discard(),
			bheUrl: url.URL{
				Scheme: "http",
				Host:   "example.com",
			},
			maxRetries: 0,
		}

		data := make(chan []any, 1)
		data <- []any{"test"}
		close(data)

		hadErrors := client.Ingest(context.Background(), data)

		require.True(t, hadErrors)
	})
}
