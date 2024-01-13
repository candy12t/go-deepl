package deepl

import (
	"bytes"
	"context"
	"encoding/json"
)

// TODO: to Functional Option Pattern
type TranslateOption struct {
	SourceLang         string   `json:"source_lang,omitempty"`
	Context            string   `json:"context,omitempty"`
	SplitSentences     string   `json:"split_sentences,omitempty"`
	PreserveFormatting bool     `json:"preserve_formatting,omitempty"`
	Formality          string   `json:"formality,omitempty"`
	GlossaryID         string   `json:"glossary_id,omitempty"`
	TagHandling        string   `json:"tag_handling,omitempty"`
	OutlineDetection   bool     `json:"outline_detection,omitempty"`
	NonSplittingTags   []string `json:"non_splitting_tags,omitempty"`
	SplittingTags      []string `json:"splitting_tags,omitempty"`
	IgnoreTags         []string `json:"ignore_tags,omitempty"`
}

type Translations struct {
	Translations []struct {
		DetectedSourceLanguage string `json:"detected_source_language"`
		Text                   string `json:"text"`
	} `json:"translations"`
}

func (c *Client) TranslateText(ctx context.Context, text []string, targetLang string, option TranslateOption) (*Translations, error) {
	params := struct {
		Text       []string `json:"text"`
		TargetLang string   `json:"target_lang"`
		TranslateOption
	}{
		Text:            text,
		TargetLang:      targetLang,
		TranslateOption: option,
	}
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	translations := new(Translations)
	if _, err := c.Post(ctx, "/translate", bytes.NewReader(body), translations); err != nil {
		return nil, err
	}
	return translations, nil
}
