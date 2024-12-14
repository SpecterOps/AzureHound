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

package cmd

import (
	"bufio"
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"io/fs"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"runtime/pprof"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/net/proxy"

	"github.com/bloodhoundad/azurehound/v2/client"
	client_config "github.com/bloodhoundad/azurehound/v2/client/config"
	"github.com/bloodhoundad/azurehound/v2/client/rest"
	"github.com/bloodhoundad/azurehound/v2/config"
	"github.com/bloodhoundad/azurehound/v2/enums"
	"github.com/bloodhoundad/azurehound/v2/logger"
	"github.com/bloodhoundad/azurehound/v2/models"
	"github.com/bloodhoundad/azurehound/v2/pipeline"
	"github.com/bloodhoundad/azurehound/v2/sinks"
)

func exit(err error) {
	log.Error(err, "encountered unrecoverable error")
	log.GetSink()
	os.Exit(1)
}

func persistentPreRunE(cmd *cobra.Command, args []string) error {
	// need to set config flag value explicitly
	if cmd != nil {
		if configFlag := cmd.Flag(config.ConfigFile.Name).Value.String(); configFlag != "" {
			config.ConfigFile.Set(configFlag)
		}
	}

	config.LoadValues(cmd, config.Options())
	config.SetAzureDefaults()

	if logr, err := logger.GetLogger(); err != nil {
		return err
	} else {
		log = *logr
		config.CheckCollectionConfigSanity(log)

		if config.ConfigFileUsed() != "" {
			log.V(1).Info(fmt.Sprintf("Config File: %v", config.ConfigFileUsed()))
		}

		if config.LogFile.Value() != "" {
			log.V(1).Info(fmt.Sprintf("Log File: %v", config.LogFile.Value()))
		}

		return nil
	}
}

func gracefulShutdown(stop context.CancelFunc) {
	stop()
	fmt.Fprintln(os.Stderr, "\nshutting down gracefully, press ctrl+c again to force")
	if profile := pprof.Lookup(config.Pprof.Value().(string)); profile != nil {
		profile.WriteTo(os.Stderr, 1)
	}
}

func testConnections() error {
	if _, err := dial(config.AzAuthUrl.Value().(string)); err != nil {
		return fmt.Errorf("unable to connect to %s: %w", config.AzAuthUrl.Value(), err)
	} else if _, err := dial(config.AzGraphUrl.Value().(string)); err != nil {
		return fmt.Errorf("unable to connect to %s: %w", config.AzGraphUrl.Value(), err)
	} else if _, err := dial(config.AzMgmtUrl.Value().(string)); err != nil {
		return fmt.Errorf("unable to connect to %s: %w", config.AzMgmtUrl.Value(), err)
	} else {
		return nil
	}
}

type httpsDialer struct{}

func (s httpsDialer) Dial(network string, addr string) (net.Conn, error) {
	return tls.Dial(network, addr, &tls.Config{})
}

func newProxyDialer(url *url.URL, forward proxy.Dialer) (proxy.Dialer, error) {
	dialer := &proxyDialer{
		host:    url.Host,
		forward: forward,
	}

	if url.User != nil {
		dialer.user = url.User.Username()
		dialer.pass, _ = url.User.Password()
	}

	return dialer, nil
}

type proxyDialer struct {
	host    string
	user    string
	pass    string
	forward proxy.Dialer
}

func (s proxyDialer) Dial(network string, addr string) (net.Conn, error) {
	if s.forward == nil {
		return nil, fmt.Errorf("unable to connect to %s: forward dialer not set", s.host)
	} else if conn, err := s.forward.Dial(network, s.host); err != nil {
		return nil, fmt.Errorf("unable to connect to %s: %w", s.host, err)
	} else if req, err := http.NewRequest("CONNECT", "//"+addr, nil); err != nil {
		conn.Close()
		return nil, fmt.Errorf("unable to connect to %s: %w", addr, err)
	} else {
		req.Close = false
		if s.user != "" {
			req.SetBasicAuth(s.user, s.pass)
		}

		// Write request over proxy connection
		if err := req.Write(conn); err != nil {
			conn.Close()
			return nil, fmt.Errorf("unable to connect to %s: %w", addr, err)
		}

		res, err := http.ReadResponse(bufio.NewReader(conn), req)
		defer func() {
			if res.Body != nil {
				res.Body.Close()
			}
		}()

		if err != nil {
			conn.Close()
			return nil, fmt.Errorf("unable to connect to %s: %w", addr, err)
		} else if res.StatusCode != 200 {
			if res.Body != nil {
				res.Body.Close()
			}
			conn.Close()
			return nil, fmt.Errorf("unable to connect to %s via proxy (%s): statusCode %d", addr, s.host, res.StatusCode)
		} else {
			return conn, nil
		}
	}
}

func getDialer() (proxy.Dialer, error) {
	if proxyUrl := config.Proxy.Value().(string); proxyUrl == "" {
		return proxy.Direct, nil
	} else if url, err := url.Parse(proxyUrl); err != nil {
		return nil, err
	} else if url.Scheme == "https" {
		return proxy.FromURL(url, httpsDialer{})
	} else {
		return proxy.FromURL(url, proxy.Direct)
	}
}

func init() {
	proxy.RegisterDialerType("http", newProxyDialer)
	proxy.RegisterDialerType("https", newProxyDialer)
}

func dial(targetUrl string) (string, error) {
	log.V(2).Info("dialing...", "targetUrl", targetUrl)
	if dialer, err := getDialer(); err != nil {
		return "", err
	} else if url, err := url.Parse(targetUrl); err != nil {
		return "", err
	} else {
		port := url.Port()

		if port == "" {
			port = "443"
		}

		if conn, err := dialer.Dial("tcp", fmt.Sprintf("%s:%s", url.Hostname(), port)); err != nil {
			return "", err
		} else {
			defer conn.Close()
			addr := conn.LocalAddr().(*net.TCPAddr)
			return addr.IP.String(), nil
		}
	}
}

func newAzureClient() (client.AzureClient, error) {
	var (
		certFile   = config.AzCert.Value()
		keyFile    = config.AzKey.Value()
		clientCert string
		clientKey  string
	)

	if file, ok := certFile.(string); ok && file != "" {
		if content, err := os.ReadFile(certFile.(string)); err != nil {
			return nil, fmt.Errorf("unable to read provided certificate: %w", err)
		} else {
			clientCert = string(content)
		}
	}

	if file, ok := keyFile.(string); ok && file != "" {
		if content, err := os.ReadFile(keyFile.(string)); err != nil {
			return nil, fmt.Errorf("unable to read provided key file: %w", err)
		} else {
			clientKey = string(content)
		}
	}

	config := client_config.Config{
		ApplicationId:   config.AzAppId.Value().(string),
		Authority:       config.AzAuthUrl.Value().(string),
		ClientSecret:    config.AzSecret.Value().(string),
		ClientCert:      clientCert,
		ClientKey:       clientKey,
		ClientKeyPass:   config.AzKeyPass.Value().(string),
		Graph:           config.AzGraphUrl.Value().(string),
		JWT:             config.JWT.Value().(string),
		Management:      config.AzMgmtUrl.Value().(string),
		MgmtGroupId:     config.AzMgmtGroupId.Value().([]string),
		Password:        config.AzPassword.Value().(string),
		ProxyUrl:        config.Proxy.Value().(string),
		RefreshToken:    config.RefreshToken.Value().(string),
		Region:          config.AzRegion.Value().(string),
		SubscriptionId:  config.AzSubId.Value().([]string),
		Tenant:          config.AzTenant.Value().(string),
		Username:        config.AzUsername.Value().(string),
		ManagedIdentity: config.AzUseManagedIdentity.Value().(bool),
	}
	return client.NewClient(config)
}

func newSigningHttpClient(signature, tokenId, token, proxyUrl string) (*http.Client, error) {
	if client, err := rest.NewHTTPClient(proxyUrl); err != nil {
		return nil, err
	} else {
		client.Transport = signingTransport{
			base:      client.Transport,
			tokenId:   tokenId,
			token:     token,
			signature: signature,
		}
		return client, nil
	}
}

type rewindableByteReader struct {
	data *bytes.Reader
}

func (s *rewindableByteReader) Read(p []byte) (int, error) {
	return s.data.Read(p)
}

func (s *rewindableByteReader) Close() error {
	return nil
}

func (s *rewindableByteReader) Rewind() (int64, error) {
	return s.data.Seek(0, io.SeekStart)
}

func discard(reader io.Reader) {
	io.Copy(io.Discard, reader)
}

type signingTransport struct {
	base      http.RoundTripper
	tokenId   string
	token     string
	signature string
}

func (s signingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// The http client may try to call RoundTrip more than once to replay the same request; in which case rewind the request
	if rbr, ok := req.Body.(*rewindableByteReader); ok {
		if _, err := rbr.Rewind(); err != nil {
			return nil, err
		}
	}

	if req.Header.Get("Signature") == "" {

		// token
		digester := hmac.New(sha256.New, []byte(s.token))

		// path
		if _, err := digester.Write([]byte(req.Method + req.URL.Path)); err != nil {
			return nil, err
		}

		// datetime
		datetime := time.Now().Format(time.RFC3339)
		digester = hmac.New(sha256.New, digester.Sum(nil))
		// hash the substring of the current datetime excluding minutes, seconds, microseconds and timezone
		if _, err := digester.Write([]byte(datetime[:13])); err != nil {
			return nil, err
		}

		// body
		digester = hmac.New(sha256.New, digester.Sum(nil))
		if req.Body != nil {
			var (
				body    = &bytes.Buffer{}
				hashBuf = make([]byte, 64*1024) // 64KB buffer, consider benchmarking and optimizing this value
				tee     = io.TeeReader(req.Body, body)
			)

			defer req.Body.Close()
			defer discard(tee)
			defer discard(body)

			for {
				numRead, err := tee.Read(hashBuf)
				if numRead > 0 {
					if _, err := digester.Write(hashBuf[:numRead]); err != nil {
						return nil, err
					}
				}

				// exit loop on EOF or error
				if err != nil {
					if err != io.EOF {
						return nil, err
					}
					break
				}
			}

			req.Body = &rewindableByteReader{data: bytes.NewReader(body.Bytes())}
		}

		signature := digester.Sum(nil)

		req.Header.Set("Authorization", fmt.Sprintf("%s %s", s.signature, s.tokenId))
		req.Header.Set("RequestDate", datetime)
		req.Header.Set("Signature", base64.StdEncoding.EncodeToString(signature))
	}
	return s.base.RoundTrip(req)
}

func contains[T comparable](collection []T, value T) bool {
	for _, item := range collection {
		if item == value {
			return true
		}
	}
	return false
}

func unique(collection []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, item := range collection {
		if _, found := keys[item]; !found {
			keys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func stat(path string) (string, fs.FileInfo, error) {
	if info, err := os.Stat(path); err == nil {
		return path, info, nil
	} else {
		p := path + ".exe"
		info, err := os.Stat(p)
		return p, info, err
	}
}

func getExePath() (string, error) {
	exe := os.Args[0]
	if exePath, err := filepath.Abs(exe); err != nil {
		return "", err
	} else if path, info, err := stat(exePath); err != nil {
		return "", err
	} else if info.IsDir() {
		return "", fmt.Errorf("%s is a directory", path)
	} else {
		return path, nil
	}
}

func setupLogger() {
	if logger, err := logger.GetLogger(); err != nil {
		panic(err)
	} else {
		log = *logger
	}
}

// deprecated: use azureWrapper instead
type AzureWrapper struct {
	Kind enums.Kind  `json:"kind"`
	Data interface{} `json:"data"`
}

type azureWrapper[T any] struct {
	Kind enums.Kind `json:"kind"`
	Data T          `json:"data"`
}

func NewAzureWrapper[T any](kind enums.Kind, data T) azureWrapper[T] {
	return azureWrapper[T]{
		Kind: kind,
		Data: data,
	}
}

func outputStream[T any](ctx context.Context, stream <-chan T) {
	formatted := pipeline.FormatJson(ctx.Done(), stream)
	if path := config.OutputFile.Value().(string); path != "" {
		if err := sinks.WriteToFile(ctx, path, formatted); err != nil {
			exit(fmt.Errorf("failed to write stream to file: %w", err))
		}
	} else {
		sinks.WriteToConsole(ctx, formatted)
	}
}

func kvRoleAssignmentFilter(roleId string) func(models.KeyVaultRoleAssignment) bool {
	return func(ra models.KeyVaultRoleAssignment) bool {
		return path.Base(ra.RoleAssignment.Properties.RoleDefinitionId) == roleId
	}
}

func vmRoleAssignmentFilter(roleId string) func(models.VirtualMachineRoleAssignment) bool {
	return func(ra models.VirtualMachineRoleAssignment) bool {
		return path.Base(ra.RoleAssignment.Properties.RoleDefinitionId) == roleId
	}
}

func rgRoleAssignmentFilter(roleId string) func(models.ResourceGroupRoleAssignment) bool {
	return func(ra models.ResourceGroupRoleAssignment) bool {
		return path.Base(ra.RoleAssignment.Properties.RoleDefinitionId) == roleId
	}
}

func mgmtGroupRoleAssignmentFilter(roleId string) func(models.ManagementGroupRoleAssignment) bool {
	return func(ra models.ManagementGroupRoleAssignment) bool {
		return path.Base(ra.RoleAssignment.Properties.RoleDefinitionId) == roleId
	}
}

func connectAndCreateClient() client.AzureClient {
	log.V(1).Info("testing connections")
	if err := testConnections(); err != nil {
		exit(fmt.Errorf("failed to test connections: %w", err))
	} else if azClient, err := newAzureClient(); err != nil {
		exit(fmt.Errorf("failed to create new Azure client: %w", err))
	} else {
		return azClient
	}

	panic("unexpectedly failed to create azClient without error")
}
