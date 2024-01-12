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
	BaseURL    *url.URL
	AuthKey    string
	UserAgent  string
}

func NewClient(authKey string) *Client {
	u, _ := url.Parse(apiHost(authKey))
	baseURL := u.JoinPath(APIVersion)
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

func (c *Client) NewRequest(ctx context.Context, method, path string, query url.Values, body io.Reader) (*http.Request, error) {
	u := c.BaseURL.JoinPath(path)
	u.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, method, u.String(), body)
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

	if v == nil {
		if _, err := io.ReadAll(resp.Body); err != nil {
			return nil, err
		}
		return resp, nil
	}

	return resp, json.NewDecoder(resp.Body).Decode(v)
}

func (c *Client) httpClient() *http.Client {
	if c.HTTPClient != nil {
		return c.HTTPClient
	}
	return http.DefaultClient
}

func (c *Client) Get(ctx context.Context, path string, query url.Values, v any) (*http.Response, error) {
	req, err := c.NewRequest(ctx, http.MethodGet, path, query, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req, v)
}

func (c *Client) Post(ctx context.Context, path string, body io.Reader, v any) (*http.Response, error) {
	req, err := c.NewRequest(ctx, http.MethodPost, path, nil, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return c.Do(req, v)
}

// TODO
func (c *Client) UploadFile(ctx context.Context, path string, body io.Reader, v any) (*http.Response, error) {
	return nil, nil
}

func (c *Client) Delete(ctx context.Context, path string) (*http.Response, error) {
	req, err := c.NewRequest(ctx, http.MethodDelete, path, nil, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req, nil)
}
