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

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/certmichelin/azurehound/v2/config"
	"github.com/certmichelin/azurehound/v2/constants"
)

func NewHTTPClient(proxyUrl string) (*http.Client, error) {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxConnsPerHost = config.ColMaxConnsPerHost.Value().(int)
	transport.MaxIdleConnsPerHost = config.ColMaxIdleConnsPerHost.Value().(int)
	transport.DisableKeepAlives = false

	// defaults to TLS 1.0 which is not favorable
	transport.TLSClientConfig = &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	// increasing timeout because tls handshakes can take longer when doing a lot of concurrent calls
	transport.TLSHandshakeTimeout = 20 * time.Second

	// increasing response header timeout to accout for WAF throttling rules
	transport.ResponseHeaderTimeout = 5 * time.Minute

	// ignoring err; always nil
	jar, _ := cookiejar.New(nil)

	// setup forward proxy
	if proxyUrl != "" {
		if url, err := url.Parse(proxyUrl); err != nil {
			return nil, err
		} else {
			transport.Proxy = http.ProxyURL(url)
		}
	}

	return &http.Client{
		Jar:       jar,
		Transport: transport,
	}, nil
}

func NewRequest(
	ctx context.Context,
	verb string,
	endpoint *url.URL,
	body interface{},
	params map[string]string,
	headers map[string]string,
) (*http.Request, error) {
	// set query params
	if params != nil {
		q := endpoint.Query()
		for key, value := range params {
			q.Set(key, value)
		}
		endpoint.RawQuery = q.Encode()
	}

	// set body
	var (
		reader io.Reader
		buffer = &bytes.Buffer{}
	)
	if body != nil {
		switch body := body.(type) {
		case url.Values:
			reader = strings.NewReader(body.Encode())
		default:
			data := new(bytes.Buffer)
			if err := json.NewEncoder(data).Encode(body); err != nil {
				return nil, err
			} else {
				reader = data
			}
		}
		buffer.ReadFrom(reader)
	}

	if req, err := http.NewRequestWithContext(ctx, verb, endpoint.String(), buffer); err != nil {
		return nil, err
	} else {
		// set headers
		for key, value := range headers {
			req.Header.Set(key, value)
		}

		// set content-type
		if body != nil {
			switch body.(type) {
			case url.Values:
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			default:
				req.Header.Set("Content-Type", "application/json")
			}
		}

		// set default accept type
		if req.Header.Get("Accept") == "" {
			req.Header.Set("Accept", "application/json")
		}

		// set azurehound as user-agent
		req.Header.Set("User-Agent", constants.UserAgent())
		return req, nil
	}
}
