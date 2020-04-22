package weather

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"weather/internal/httpclient"
)

func TestClient_GetByCity_Errors(t *testing.T) {
	tcs := []struct {
		n string
		c *Client
	}{
		{
			n: "bad url",
			c: &Client{baseURL: "!@#$%"},
		},
		{
			n: "failed request",
			c: &Client{baseURL: "http://localhost", d: httpclient.DoerFunc(func(*http.Request) (*http.Response, error) {
				return nil, errors.New("error")
			})},
		}, {
			n: "bad request body",
			c: &Client{baseURL: "http://localhost", d: httpclient.DoerFunc(func(*http.Request) (*http.Response, error) {
				return &http.Response{
					Body: failingReader(true),
				}, nil
			})},
		}, {
			n: "bad code",
			c: &Client{baseURL: "http://localhost", d: httpclient.DoerFunc(func(*http.Request) (*http.Response, error) {
				return &http.Response{
					Body: ioutil.NopCloser(strings.NewReader("d")),
				}, nil
			})},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.n, func(t *testing.T) {
			_, err := tc.c.GetByCity("")
			if err == nil {
				t.Error("expected an error")
			}
		})
	}
}

type failingReader bool

func (r failingReader) err(m string) error {
	if r {
		return errors.New("FAIL " + m)
	}

	return nil
}
func (r failingReader) Close() error             { return nil }
func (r failingReader) Read([]byte) (int, error) { return 0, r.err("Read") }
