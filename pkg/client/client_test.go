package client_test

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/flatcar-linux/fleetlock/pkg/client"
)

type mockRoundTripper struct {
	resp  *http.Response
	r     *http.Request
	doErr error
}

func (r *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	r.r = req
	return r.resp, r.doErr
}

func TestBadURL(t *testing.T) {
	_, err := client.New("this is not an URL", "default", "1234", nil)
	if err == nil {
		t.Fatalf("should get an error")
	}

	if err != nil && err.Error() != "parsing URL: parse \"this is not an URL\": invalid URI for request" {
		t.Fatalf("should have %s for the error, got: %s", "parsing URL: parse \"this is not an URL\": invalid URI for request", err.Error())
	}
}

func TestRecursiveLock(t *testing.T) {

	expURL := "http://1.2.3.4/v1/pre-reboot"

	for _, test := range []struct {
		statusCode int
		expErr     error
		body       io.ReadCloser
		doErr      error
	}{
		{
			statusCode: 200,
		},
		{
			statusCode: 500,
			expErr:     errors.New("fleetlock error: this is an error (error_kind)"),
			body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"kind":"error_kind","value":"this is an error"}`))),
		},
		{
			statusCode: 500,
			expErr:     errors.New("unmarshalling error: invalid character '\"' after object key:value pair"),
			body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"kind":"error_kind"  "value":"this is an error"}`))),
		},
		{
			statusCode: 100,
			expErr:     errors.New("unexpected status code: 100"),
		},
		{
			expErr: errors.New("doing the request: Post \"http://1.2.3.4/v1/pre-reboot\": connection refused"),
			doErr:  errors.New("connection refused"),
		},
	} {
		h := http.DefaultClient
		tr := &mockRoundTripper{
			resp: &http.Response{
				StatusCode: test.statusCode,
				Body:       test.body,
			},
			doErr: test.doErr,
		}

		h.Transport = tr

		c, _ := client.New("http://1.2.3.4", "default", "1234", h)

		err := c.RecursiveLock()
		if err != nil && err.Error() != test.expErr.Error() {
			t.Fatalf("should have %v for err, got: %v", test.expErr, err)
		}

		if tr.r.URL.String() != expURL {
			t.Fatalf("should have %s for URL, got: %s", expURL, tr.r.URL.String())
		}
	}
}

func TestUnlockIfHeld(t *testing.T) {

	expURL := "http://1.2.3.4/v1/steady-state"

	for _, test := range []struct {
		statusCode int
		expErr     error
		body       io.ReadCloser
		doErr      error
	}{
		{
			statusCode: 200,
		},
		{
			statusCode: 500,
			expErr:     errors.New("fleetlock error: this is an error (error_kind)"),
			body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"kind":"error_kind","value":"this is an error"}`))),
		},
		{
			statusCode: 500,
			expErr:     errors.New("unmarshalling error: invalid character '\"' after object key:value pair"),
			body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"kind":"error_kind"  "value":"this is an error"}`))),
		},
		{
			statusCode: 100,
			expErr:     errors.New("unexpected status code: 100"),
		},
		{
			expErr: errors.New("doing the request: Post \"http://1.2.3.4/v1/steady-state\": connection refused"),
			doErr:  errors.New("connection refused"),
		},
	} {
		h := http.DefaultClient
		tr := &mockRoundTripper{
			resp: &http.Response{
				StatusCode: test.statusCode,
				Body:       test.body,
			},
			doErr: test.doErr,
		}

		h.Transport = tr

		c, _ := client.New("http://1.2.3.4", "default", "1234", h)

		err := c.UnlockIfHeld()
		if err != nil && err.Error() != test.expErr.Error() {
			t.Fatalf("should have %v for err, got: %v", test.expErr, err)
		}

		if tr.r.URL.String() != expURL {
			t.Fatalf("should have %s for URL, got: %s", expURL, tr.r.URL.String())
		}
	}
}
