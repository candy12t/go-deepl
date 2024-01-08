package deepl

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	Version = "v0.0.1"

	ProAPIHost  = "https://api.deepl.com"
	FreeAPIHost = "https://api-free.deepl.com"
	APIVersion  = "v2"

	defaultUserAgent = "go-deepl" + "/" + Version + " (https://https://github.com/candy12t/go-deepl)"

	accountPlanIdentifyKey = ":fx"
)

type Client struct {
	HTTPClient *http.Client
	BaseURL    string
	AuthKey    string
	UserAgent  string
}

func NewClient(authKey string) *Client {
	baseURL, _ := url.JoinPath(apiHost(authKey), APIVersion)
	return &Client{
		BaseURL:   baseURL,
		AuthKey:   authKey,
		UserAgent: defaultUserAgent,
	}
}

func isFreeAccount(authKey string) bool {
	return strings.HasSuffix(authKey, accountPlanIdentifyKey)
}

func apiHost(authKey string) string {
	if isFreeAccount(authKey) {
		return FreeAPIHost
	}
	return ProAPIHost
}

func (c *Client) NewRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	u, err := url.JoinPath(c.BaseURL, path)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, method, u, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("DeepL-Auth-Key %s", c.AuthKey))
	req.Header.Set("User-Agent", c.UserAgent)

	return req, nil
}

func (c *Client) Do(req *http.Request, v any) (*http.Response, error) {
	resp, err := c.httpClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	ok := http.StatusOK <= resp.StatusCode && resp.StatusCode < http.StatusMultipleChoices
	if !ok {
		return nil, HandleHTTPError(resp)
	}

	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) httpClient() *http.Client {
	if c.HTTPClient != nil {
		return c.HTTPClient
	}
	return http.DefaultClient
}
