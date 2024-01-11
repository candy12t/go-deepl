package deepl

import (
	"context"
	"net/http"
)

type Usage struct {
	CharacterCount    int `json:"character_count"`
	CharacterLimit    int `json:"character_limit"`
	DocumentCount     int `json:"document_count,omitempty"`
	DocumentLimit     int `json:"document_limit,omitempty"`
	TeamDocumentCount int `json:"team_document_count,omitempty"`
	TeamDocumentLimit int `json:"team_document_limit,omitempty"`
}

func (c *Client) GetUsage(ctx context.Context) (*Usage, error) {
	req, err := c.NewRequest(ctx, http.MethodGet, "/usage", nil, nil)
	if err != nil {
		return nil, err
	}

	usage := new(Usage)
	if _, err := c.Do(req, usage); err != nil {
		return nil, err
	}
	return usage, nil
}
