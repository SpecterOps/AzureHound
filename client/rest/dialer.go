package rest

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"

	"github.com/bloodhoundad/azurehound/v2/config"
	"github.com/go-logr/logr"
	"golang.org/x/net/proxy"
)

type HttpsDialer struct{}

func (s HttpsDialer) Dial(network string, addr string) (net.Conn, error) {
	// TODO: look into pinning this to TLS 1.2 to avoid downgraded TLS attacks
	return tls.Dial(network, addr, &tls.Config{})
}

func NewProxyDialer(url *url.URL, forward proxy.Dialer) (proxy.Dialer, error) {
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

func GetDialer() (proxy.Dialer, error) {
	if proxyUrl := config.Proxy.Value().(string); proxyUrl == "" {
		return proxy.Direct, nil
	} else if url, err := url.Parse(proxyUrl); err != nil {
		return nil, err
	} else if url.Scheme == "https" {
		return proxy.FromURL(url, HttpsDialer{})
	} else {
		return proxy.FromURL(url, proxy.Direct)
	}
}

func Dial(log logr.Logger, targetUrl string) (string, error) {
	log.V(2).Info("dialing...", "targetUrl", targetUrl)
	if dialer, err := GetDialer(); err != nil {
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
