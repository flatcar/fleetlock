package client_test

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/flatcar-linux/fleetlock/pkg/client"
)

type httpClient struct {
	resp  *http.Response
	r     *http.Request
	doErr error
}

func (h *httpClient) Do(req *http.Request) (*http.Response, error) {
	h.r = req

	return h.resp, h.doErr
}

func TestBadURL(t *testing.T) {
	t.Parallel()

	_, err := client.New("this is not an URL", "default", "1234", nil)
	if err == nil {
		t.Fatalf("should get an error")
	}

	expectedError := "parsing URL: parse \"this is not an URL\": invalid URI for request"
	if err != nil && err.Error() != expectedError {
		t.Fatalf("should have %s for the error, got: %s", expectedError, err.Error())
	}
}

//nolint:funlen // Just many test cases.
func TestClient(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		statusCode int
		expErr     error
		body       []byte
		doErr      error
	}{
		{
			statusCode: 200,
		},
		{
			statusCode: 500,
			expErr:     errors.New("fleetlock error: this is an error (error_kind)"),
			body:       []byte(`{"kind":"error_kind","value":"this is an error"}`),
		},
		{
			statusCode: 500,
			expErr:     errors.New("unmarshalling error: invalid character '\"' after object key:value pair"),
			body:       []byte(`{"kind":"error_kind"  "value":"this is an error"}`),
		},
		{
			statusCode: 100,
			expErr:     errors.New("unexpected status code: 100"),
		},
		{
			expErr: errors.New("doing the request: connection refused"),
			doErr:  errors.New("connection refused"),
		},
	} {
		test := test

		newClient := func(statusCode int, body []byte, doErr error) (*httpClient, *client.Client) {
			h := &httpClient{
				resp: &http.Response{
					StatusCode: statusCode,
					Body:       ioutil.NopCloser(bytes.NewReader(body)),
				},
				doErr: doErr,
			}

			c, _ := client.New("http://1.2.3.4", "default", "1234", h)

			return h, c
		}

		t.Run(fmt.Sprintf("UnlockIfHeld_%d", test.statusCode), func(t *testing.T) {
			t.Parallel()

			h, c := newClient(test.statusCode, test.body, test.doErr)

			err := c.UnlockIfHeld()
			if err != nil && err.Error() != test.expErr.Error() {
				t.Fatalf("should have %v for err, got: %v", test.expErr, err)
			}

			expURL := "http://1.2.3.4/v1/steady-state"

			if h.r.URL.String() != expURL {
				t.Fatalf("should have %s for URL, got: %s", expURL, h.r.URL.String())
			}
		})

		t.Run(fmt.Sprintf("RecursiveLock_%d", test.statusCode), func(t *testing.T) {
			t.Parallel()

			h, c := newClient(test.statusCode, test.body, test.doErr)

			err := c.RecursiveLock()
			if err != nil && err.Error() != test.expErr.Error() {
				t.Fatalf("should have %v for err, got: %v", test.expErr, err)
			}

			expURL := "http://1.2.3.4/v1/pre-reboot"

			if h.r.URL.String() != expURL {
				t.Fatalf("should have %s for URL, got: %s", expURL, h.r.URL.String())
			}
		})
	}
}
