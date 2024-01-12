package deepl

import (
	"context"
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
	var languages []Language

	query := url.Values{}
	query.Add("type", string(langType))

	if _, err := c.Get(ctx, "/languages", query, &languages); err != nil {
		return nil, err
	}
	return languages, nil
}

type GlossaryLanguagePairs struct {
	SupportedLanguages []struct {
		SourceLang string `json:"source_lang"`
		TargetLang string `json:"target_lang"`
	} `json:"supported_languages"`
}

func (c *Client) GetGlossaryLanguagesPairs(ctx context.Context) (*GlossaryLanguagePairs, error) {
	languages := new(GlossaryLanguagePairs)
	if _, err := c.Get(ctx, "/glossary-language-pairs", nil, &languages); err != nil {
		return nil, err
	}
	return languages, nil
}
