package bloodhound

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"time"
)

type roundTripperFunc func(req *http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
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

func discard(reader io.Reader) {
	io.Copy(io.Discard, reader)
}
