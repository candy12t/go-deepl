package deepl

import (
	"context"
	"net/http"
	"net/url"
)

type langType string

const (
	Source langType = "source"
	Target langType = "target"
)

type Language struct {
	Language          string `json:"language"`
	Name              string `json:"name"`
	SupportsFormality *bool  `json:"supports_formality,omitempty"`
}

func (c *Client) GetLanguages(ctx context.Context, langType langType) ([]Language, error) {
	query := url.Values{}
	query.Add("type", string(langType))

	req, err := c.NewRequest(ctx, http.MethodGet, "/languages", query, nil)
	if err != nil {
		return nil, err
	}

	var languages []Language
	if _, err := c.Do(req, &languages); err != nil {
		return nil, err
	}
	return languages, nil
}
