package main_test

import (
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestAuthoriseBadClientId(t *testing.T) {
	testCases := []struct {
		name string
		url  url.URL
	}{
		{"Missing client ID", url.URL{
			Scheme: "http",
			Host:   net.JoinHostPort("localhost", "3000"),
			Path:   "authorise",
		}},
		{"Wrong client ID", url.URL{
			Scheme:   "http",
			Host:     net.JoinHostPort("localhost", "3000"),
			Path:     "authorise",
			RawQuery: "client_id=wrong-client-id",
		}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := http.Get(tc.url.String())

			if err != nil {
				t.Errorf("Error getting user links: %v", err)
			}

			if res.StatusCode != http.StatusBadRequest {
				t.Errorf("Wrong status code: %v", res.StatusCode)
			}

			body, err := io.ReadAll(res.Body)
			res.Body.Close()

			if strings.Compare(string(body), "Unknown client\n") != 0 {
				t.Errorf("Wrong response body: %q", string(body))
			}
		})
	}
}

func TestAuthoriseBadClientRedirectUri(t *testing.T) {
	u := url.URL{
		Scheme:   "http",
		Host:     net.JoinHostPort("localhost", "3000"),
		Path:     "authorise",
		RawQuery: "client_id=tmp-oauth-client&redirect_uri=http://localhost:8000/wrong-oauth-redirect-uri",
	}

	res, err := http.Get(u.String())

	if err != nil {
		t.Errorf("Error getting user links: %v", err)
	}

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("Wrong status code: %v", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	res.Body.Close()

	if strings.Compare(string(body), "Invalid redirect URI\n") != 0 {
		t.Errorf("Wrong response body: %q", string(body))
	}
}

func TestAuthoriseSuccessful(t *testing.T) {
	u := url.URL{
		Scheme:   "http",
		Host:     net.JoinHostPort("localhost", "3000"),
		Path:     "authorise",
		RawQuery: "client_id=tmp-oauth-client&redirect_uri=http://localhost:8000/tmp-oauth-redirect",
	}

	res, err := http.Get(u.String())

	if err != nil {
		t.Errorf("Error getting user links: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("Wrong status code: %v", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	res.Body.Close()

	if len(body) < 26 {
		t.Errorf("Wrong random key length %d", len(body))
	}
}
