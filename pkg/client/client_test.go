package client_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/flatcar/fleetlock/pkg/client"
)

type httpClient struct {
	do func(req *http.Request) (*http.Response, error)
	r  *http.Request
}

func (m *httpClient) Do(req *http.Request) (*http.Response, error) {
	m.r = req

	return m.do(req)
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

	ctx := context.Background()

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
				do: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: statusCode,
						Body:       ioutil.NopCloser(bytes.NewReader(body)),
					}, doErr
				},
			}

			c, err := client.New("http://1.2.3.4", "default", "1234", h)
			if err != nil {
				t.Fatalf("Unexpected error creating client: %v", err)
			}

			return h, c
		}

		t.Run(fmt.Sprintf("UnlockIfHeld_%d", test.statusCode), func(t *testing.T) {
			t.Parallel()

			h, c := newClient(test.statusCode, test.body, test.doErr)

			err := c.UnlockIfHeld(ctx)
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

			err := c.RecursiveLock(ctx)
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

func Test_Client_use_given_context_for_requests(t *testing.T) {
	t.Parallel()

	key := struct{}{}
	value := "bar"

	h := &httpClient{
		do: func(req *http.Request) (*http.Response, error) {
			if req.Context().Value(key) == nil {
				t.Fatalf("Expected request to use given context")
			}

			return &http.Response{
				StatusCode: 200,
			}, nil
		},
	}

	c, err := client.New("http://1.2.3.4", "default", "1234", h)
	if err != nil {
		t.Fatalf("Unexpected error creating client: %v", err)
	}

	ctx := context.WithValue(context.Background(), key, value)

	if err := c.RecursiveLock(ctx); err != nil {
		t.Fatalf("Unexpected error while doing recursive lock: %v", err)
	}

	if err := c.UnlockIfHeld(ctx); err != nil {
		t.Fatalf("Unexpected error while unlocking: %v", err)
	}
}
