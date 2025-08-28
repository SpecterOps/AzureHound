// Copyright (C) 2024 Specter Ops, Inc.
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
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/bloodhoundad/azurehound/v2/client/config"
	"github.com/bloodhoundad/azurehound/v2/constants"
)

// AuthStrategy is an interface that defines the methods that an authentication strategy must implement
type AuthStrategy interface {
	isExpired() bool
	createAuthRequest() (*http.Request, error)
	decodeAuthResponse(resp *http.Response) error
	addAuthenticationToRequest(req *http.Request) (*http.Request, error)
}

// Authenticator manages the authentication process, using a specific AuthStrategy
type Authenticator struct {
	auth  AuthStrategy
	mutex sync.RWMutex
}

// ManagedIdentityAuthStrategy is an authentication strategy that uses Azure Managed Identity
type ManagedIdentityAuthStrategy struct {
	config  config.Config
	authUrl url.URL
	api     url.URL
	tenant  string
	token   Token
}

// ManagedIdentitySDKAuthStrategy is an authentication strategy that uses Azure Managed Identity SDK
type ManagedIdentitySDKAuthStrategy struct {
	credential *azidentity.ManagedIdentityCredential
	resource   url.URL
	token      azcore.AccessToken
}

// GenericAuthStrategy is an authentication strategy that uses a bunch of pre-existing authentication methods (TODO: Break this up)
type GenericAuthStrategy struct {
	config        config.Config
	api           url.URL
	authUrl       url.URL
	jwt           string
	clientId      string
	clientSecret  string
	clientCert    string
	clientKey     string
	clientKeyPass string
	username      string
	password      string
	refreshToken  string
	tenant        string
	token         Token
}

// NewManagedIdentityAuthenticator creates a new Authenticator using the ManagedIdentityAuthStrategy
func NewManagedIdentityAuthenticator(config config.Config, auth *url.URL, api *url.URL, http *http.Client) *Authenticator {
	return &Authenticator{
		auth: &ManagedIdentityAuthStrategy{
			config:  config,
			authUrl: *auth,
			api:     *api,
			tenant:  config.Tenant,
		},
		mutex: sync.RWMutex{},
	}
}

// NewManagedIdentitySDKAuthenticator creates a new Authenticator using the ManagedIdentitySDKAuthStrategy
func NewManagedIdentitySDKAuthenticator(config config.Config, resource *url.URL) *Authenticator {

	options := &azidentity.ManagedIdentityCredentialOptions{
		ID: azidentity.ClientID(config.ManagedIdentityClientId),
	}

	cred, err := azidentity.NewManagedIdentityCredential(options)
	if err != nil {
		log.Fatalf("Failed to authenticate with User-Assigned Managed Identity")
		optionsFallback := &azidentity.ManagedIdentityCredentialOptions{}
		cred, err = azidentity.NewManagedIdentityCredential(optionsFallback)
		if err != nil {
			log.Fatalf("Failed to authenticate with system-assigned Managed Identity: %v", err)
			return nil
		}
	}
	return &Authenticator{
		auth: &ManagedIdentitySDKAuthStrategy{
			credential: cred,
			resource:   *resource,
		},
		mutex: sync.RWMutex{},
	}
}

// NewGenericAuthenticator creates a new Authenticator using the GenericAuthStrategy (The collection of pre-existing authentication methods)
func NewGenericAuthenticator(config config.Config, auth *url.URL, api *url.URL) *Authenticator {
	return &Authenticator{
		auth: &GenericAuthStrategy{config: config,
			authUrl:       *auth,
			api:           *api,
			jwt:           config.JWT,
			clientId:      config.ApplicationId,
			clientSecret:  config.ClientSecret,
			clientCert:    config.ClientCert,
			clientKey:     config.ClientKey,
			clientKeyPass: config.ClientKeyPass,
			username:      config.Username,
			password:      config.Password,
			refreshToken:  config.RefreshToken,
			tenant:        config.Tenant,
			token:         Token{},
		},
		mutex: sync.RWMutex{},
	}
}

// Authenticate if needed and add authentication to the request
func (s *Authenticator) AddAuthenticationToRequest(restClient *restClient, req *http.Request) (*http.Request, error) {
	if err := s.refreshIfExpired(restClient); err != nil {
		return nil, err
	}
	if req, err := s.auth.addAuthenticationToRequest(req); err != nil {
		return nil, err
	} else {
		return req, err
	}
}

// Authenticate if needed using a specific AuthStrategy
func (s *Authenticator) refreshIfExpired(r *restClient) error {
	if !s.auth.isExpired() {
		return nil
	}
	// Authenticate
	authRequest, err := s.auth.createAuthRequest()
	if err != nil {
		return err
	}
	if authRequest == nil {
		return nil
	}
	if authResponse, err := r.send(authRequest); err != nil {
		return err
	} else {
		defer authResponse.Body.Close()
		s.mutex.Lock()
		defer s.mutex.Unlock()

		if err := s.auth.decodeAuthResponse(authResponse); err != nil {
			return err
		}
	}
	return nil
}

func (s *ManagedIdentityAuthStrategy) isExpired() bool {
	return s.token.IsExpired()
}

func (s *ManagedIdentityAuthStrategy) addAuthenticationToRequest(req *http.Request) (*http.Request, error) {
	req.Header.Set("Authorization", s.token.String())

	return req, nil
}

func (s *ManagedIdentityAuthStrategy) createAuthRequest() (*http.Request, error) {
	endpoint, err := url.Parse("http://169.254.169.254/metadata/identity/oauth2/token?api-version=2018-02-01")
	if err != nil {
		return nil, err
	}

	getArgs := endpoint.Query()
	getArgs.Add("resource", s.api.String())
	endpoint.RawQuery = getArgs.Encode()

	req, err := NewRequest(context.Background(), "GET", endpoint, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Metadata", "true")

	return req, nil
}

func (s *ManagedIdentityAuthStrategy) decodeAuthResponse(resp *http.Response) error {
	if err := json.NewDecoder(resp.Body).Decode(&s.token); err != nil {
		return err
	} else {
		return nil
	}
}

func (s *GenericAuthStrategy) createAuthRequest() (*http.Request, error) {
	var (
		path         = url.URL{Path: fmt.Sprintf("/%s/oauth2/v2.0/token", s.tenant)}
		endpoint     = s.authUrl.ResolveReference(&path)
		defaultScope = url.URL{Path: "/.default"}
		scope        = s.api.ResolveReference(&defaultScope)
		body         = url.Values{}
	)

	if s.clientId == "" {
		body.Add("client_id", constants.AzPowerShellClientID)
	} else {
		body.Add("client_id", s.clientId)
	}

	body.Add("scope", scope.ResolveReference(&defaultScope).String())

	if s.refreshToken != "" {
		body.Add("grant_type", "refresh_token")
		body.Add("refresh_token", s.refreshToken)
		body.Set("client_id", constants.AzPowerShellClientID)
	} else if s.clientSecret != "" {
		body.Add("grant_type", "client_credentials")
		body.Add("client_secret", s.clientSecret)
	} else if s.clientCert != "" && s.clientKey != "" {
		if clientAssertion, err := NewClientAssertion(endpoint.String(), s.clientId, s.clientCert, s.clientKey, s.clientKeyPass); err != nil {
			return nil, err
		} else {
			body.Add("grant_type", "client_credentials")
			body.Add("client_assertion_type", "urn:ietf:params:oauth:client-assertion-type:jwt-bearer")
			body.Add("client_assertion", clientAssertion)
		}
	} else if s.username != "" && s.password != "" {
		body.Add("grant_type", "password")
		body.Add("username", s.username)
		body.Add("password", s.password)
		body.Set("client_id", constants.AzPowerShellClientID)
	} else {
		return nil, fmt.Errorf("unable to authenticate. no valid credential provided")
	}

	if authRequest, err := NewRequest(context.Background(), "POST", endpoint, body, nil, nil); err != nil {
		return nil, err
	} else {
		return authRequest, nil
	}
}

func (s *GenericAuthStrategy) isExpired() bool {
	return s.token.IsExpired()
}

func (s *GenericAuthStrategy) decodeAuthResponse(resp *http.Response) error {
	if err := json.NewDecoder(resp.Body).Decode(&s.token); err != nil {
		return err
	} else {
		return nil
	}
}

func (s *GenericAuthStrategy) addAuthenticationToRequest(req *http.Request) (*http.Request, error) {
	if s.jwt != "" {
		if aud, err := ParseAud(s.jwt); err != nil {
			return nil, err
		} else if aud != s.api.String() {
			return nil, fmt.Errorf("invalid audience")
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.jwt))
	} else {
		req.Header.Set("Authorization", s.token.String())
	}
	return req, nil
}

// Adds the access token to the outgoing HTTP request
func (s *ManagedIdentitySDKAuthStrategy) addAuthenticationToRequest(req *http.Request) (*http.Request, error) {
	token := s.token

	if token.Token == "" || token.ExpiresOn.Before(time.Now().Add(2*time.Minute)) {
		if err := s.refreshToken(req.Context()); err != nil {
			return nil, err
		}
		token = s.token
	}
	req.Header.Set("Authorization", "Bearer "+token.Token)
	return req, nil
}

func (s *ManagedIdentitySDKAuthStrategy) refreshToken(ctx context.Context) error {
	token, err := s.credential.GetToken(ctx, policy.TokenRequestOptions{
		Scopes: []string{s.resource.String() + "/.default"},
	})
	if err != nil {
		return err
	}
	s.token = token
	return nil
}

func (s *ManagedIdentitySDKAuthStrategy) createAuthRequest() (*http.Request, error) {
	// Not used in SDK-based flow, return nil
	return nil, nil
}

func (s *ManagedIdentitySDKAuthStrategy) decodeAuthResponse(resp *http.Response) error {
	// Not used in SDK-based flow, do nothing
	return nil
}

func (s *ManagedIdentitySDKAuthStrategy) isExpired() bool {
	return s.token.Token == "" || s.token.ExpiresOn.Before(time.Now().Add(2*time.Minute))
}
