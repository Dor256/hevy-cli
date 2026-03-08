package middleware

import (
	"fmt"
	"github.com/spf13/viper"
	"net/http"
)

var ErrUnauthenticated = fmt.Errorf("UNAUTHENTICATED: Hevy API key is missing.")

type AuthTransport struct {
	Base http.RoundTripper
}

func (transport *AuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	apiKey := viper.GetString("api_key")
	if apiKey == "" {
		return nil, ErrUnauthenticated
	}
	req.Header.Set("api-key", apiKey)
	req.Header.Set("Content-Type", "application/json")

	return transport.Base.RoundTrip(req)
}
