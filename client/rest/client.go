// Copyright (C) 2022 Specter Ops, Inc.
//
// This file is part of AzureHound.
//
// AzureHound is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// AzureHound is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package rest

//go:generate go run go.uber.org/mock/mockgen -destination=./mocks/client.go -package=mocks . RestClient

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/bloodhoundad/azurehound/v2/client/config"
	"github.com/bloodhoundad/azurehound/v2/client/query"
)

type RestClient interface {
	Delete(ctx context.Context, path string, body interface{}, params query.Params, headers map[string]string) (*http.Response, error)
	Get(ctx context.Context, path string, params query.Params, headers map[string]string) (*http.Response, error)
	Patch(ctx context.Context, path string, body interface{}, params query.Params, headers map[string]string) (*http.Response, error)
	Post(ctx context.Context, path string, body interface{}, params query.Params, headers map[string]string) (*http.Response, error)
	Put(ctx context.Context, path string, body interface{}, params query.Params, headers map[string]string) (*http.Response, error)
	Send(req *http.Request) (*http.Response, error)
	AddAuthenticationToRequest(req *http.Request) (*http.Request, error)
	CloseIdleConnections()
}

func NewRestClient(apiUrl string, config config.Config) (RestClient, error) {

	if auth, err := url.Parse(config.AuthorityUrl()); err != nil {
		return nil, err
	} else if api, err := url.Parse(apiUrl); err != nil {
		return nil, err
	} else if http, err := NewHTTPClient(config.ProxyUrl); err != nil {
		return nil, err
	} else {
		var authenticator *Authenticator
		if config.ManagedIdentity {
			authenticator = NewManagedIdentityAuthenticator(config, auth, api, http)
		} else {
			authenticator = NewGenericAuthenticator(config, auth, api)
		}
		client := &restClient{
			*api,
			http,
			config.Tenant,
			Token{},
			config.SubscriptionId,
			config.MgmtGroupId,
			authenticator,
		}
		return client, nil
	}
}

type restClient struct {
	api           url.URL
	http          *http.Client
	tenant        string
	token         Token
	subId         []string
	mgmtGroupId   []string
	Authenticator *Authenticator
}

func (s *restClient) Delete(ctx context.Context, path string, body interface{}, params query.Params, headers map[string]string) (*http.Response, error) {
	endpoint := s.api.ResolveReference(&url.URL{Path: path})
	paramsMap := make(map[string]string)
	if params != nil {
		paramsMap = params.AsMap()
	}
	if req, err := NewRequest(ctx, http.MethodDelete, endpoint, body, paramsMap, headers); err != nil {
		return nil, err
	} else {
		return s.Send(req)
	}
}

func (s *restClient) Get(ctx context.Context, path string, params query.Params, headers map[string]string) (*http.Response, error) {
	endpoint := s.api.ResolveReference(&url.URL{Path: path})
	paramsMap := make(map[string]string)

	if params != nil {
		paramsMap = params.AsMap()
		if params.NeedsEventualConsistencyHeaderFlag() {
			if headers == nil {
				headers = make(map[string]string)
			}
			headers["ConsistencyLevel"] = "eventual"
		}
	}

	if req, err := NewRequest(ctx, http.MethodGet, endpoint, nil, paramsMap, headers); err != nil {
		return nil, err
	} else {
		return s.Send(req)
	}
}

func (s *restClient) Patch(ctx context.Context, path string, body interface{}, params query.Params, headers map[string]string) (*http.Response, error) {
	endpoint := s.api.ResolveReference(&url.URL{Path: path})
	paramsMap := make(map[string]string)
	if params != nil {
		paramsMap = params.AsMap()
	}
	if req, err := NewRequest(ctx, http.MethodPatch, endpoint, body, paramsMap, headers); err != nil {
		return nil, err
	} else {
		return s.Send(req)
	}
}

func (s *restClient) Post(ctx context.Context, path string, body interface{}, params query.Params, headers map[string]string) (*http.Response, error) {
	endpoint := s.api.ResolveReference(&url.URL{Path: path})
	paramsMap := make(map[string]string)
	if params != nil {
		paramsMap = params.AsMap()
	}
	if req, err := NewRequest(ctx, http.MethodPost, endpoint, body, paramsMap, headers); err != nil {
		return nil, err
	} else {
		return s.Send(req)
	}
}

func (s *restClient) Put(ctx context.Context, path string, body interface{}, params query.Params, headers map[string]string) (*http.Response, error) {
	endpoint := s.api.ResolveReference(&url.URL{Path: path})
	paramsMap := make(map[string]string)
	if params != nil {
		paramsMap = params.AsMap()
	}
	if req, err := NewRequest(ctx, http.MethodPost, endpoint, body, paramsMap, headers); err != nil {
		return nil, err
	} else {
		return s.Send(req)
	}
}

func (s *restClient) AddAuthenticationToRequest(req *http.Request) (*http.Request, error) {
	return s.Authenticator.AddAuthenticationToRequest(s, req)
}

func (s *restClient) Send(req *http.Request) (*http.Response, error) {
	_, err := s.AddAuthenticationToRequest(req)
	if err != nil {
		return nil, err
	}
	return s.send(req)
}

func (s *restClient) send(req *http.Request) (*http.Response, error) {
	// copy the bytes in case we need to retry the request
	if body, err := CopyBody(req); err != nil {
		return nil, err
	} else {
		var (
			res        *http.Response
			err        error
			maxRetries = 3
		)
		// Try the request up to a set number of times
		for retry := 0; retry < maxRetries; retry++ {

			// Reusing http.Request requires rewinding the request body
			// back to a working state
			if body != nil && retry > 0 {
				req.Body = io.NopCloser(bytes.NewBuffer(body))
			}

			// Try the request
			if res, err = s.http.Do(req); err != nil {
				if IsClosedConnectionErr(err) {
					fmt.Printf("remote host force closed connection while requesting %s; attempt %d/%d; trying again\n", req.URL, retry+1, maxRetries)
					ExponentialBackoff(retry)
					continue
				}
				return nil, err
			} else if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
				// Error response code handling
				// See official Retry guidance (https://learn.microsoft.com/en-us/azure/architecture/best-practices/retry-service-specific#retry-usage-guidance)
				if res.StatusCode == http.StatusTooManyRequests {
					retryAfterHeader := res.Header.Get("Retry-After")
					if retryAfter, err := strconv.ParseInt(retryAfterHeader, 10, 64); err != nil {
						return nil, fmt.Errorf("attempting to handle 429 but unable to parse retry-after header: %w", err)
					} else {
						// Wait the time indicated in the retry-after header
						time.Sleep(time.Second * time.Duration(retryAfter))
						continue
					}
				} else if res.StatusCode >= http.StatusInternalServerError {
					// Wait the time calculated by the 5 second exponential backoff
					ExponentialBackoff(retry)
					continue
				} else {
					// Not a status code that warrants a retry
					var errRes map[string]interface{}
					if err := Decode(res.Body, &errRes); err != nil {
						return nil, fmt.Errorf("malformed error response, status code: %d", res.StatusCode)
					} else {
						return nil, fmt.Errorf("%v", errRes)
					}
				}
			} else {
				// Response OK
				return res, nil
			}
		}
		return nil, fmt.Errorf("unable to complete the request after %d attempts: %w", maxRetries, err)
	}
}

func (s *restClient) CloseIdleConnections() {
	s.http.CloseIdleConnections()
}
