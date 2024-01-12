package deepl

import (
	"context"
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
	usage := new(Usage)
	if _, err := c.Get(ctx, "/usage", nil, usage); err != nil {
		return nil, err
	}
	return usage, nil
}
