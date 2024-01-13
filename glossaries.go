package deepl

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Glossaries struct {
	Glossaries []Glossary `json:"glossaries"`
}

type Glossary struct {
	GlossaryID   string    `json:"glossary_id"`
	Name         string    `json:"name"`
	Ready        bool      `json:"ready"`
	SourceLang   string    `json:"source_lang"`
	TargetLang   string    `json:"target_lang"`
	CreationTime time.Time `json:"creation_time"`
	EntryCount   int       `json:"entry_count"`
}

func (c *Client) CreateGlossary(ctx context.Context, name, sourceLang, targetLang, entries string, entriesFormat EntriesFormat) (*Glossary, error) {
	params := struct {
		Name          string `json:"name"`
		SourceLang    string `json:"source_lang"`
		TargetLang    string `json:"target_lang"`
		Entries       string `json:"entries"`
		EntriesFormat string `json:"entries_format"`
	}{
		Name:          name,
		SourceLang:    sourceLang,
		TargetLang:    targetLang,
		Entries:       entries,
		EntriesFormat: entriesFormat.String(),
	}
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	glossary := new(Glossary)
	if _, err := c.Post(ctx, "/glossaries", bytes.NewReader(body), glossary); err != nil {
		return nil, err
	}
	return glossary, nil
}

func (c *Client) GetGlossaries(ctx context.Context) (*Glossaries, error) {
	glossaries := new(Glossaries)
	if _, err := c.Get(ctx, "/glossaries", nil, glossaries); err != nil {
		return nil, err
	}
	return glossaries, nil
}

func (c *Client) GetGlossary(ctx context.Context, glossaryID string) (*Glossary, error) {
	glossary := new(Glossary)
	if _, err := c.Get(ctx, fmt.Sprintf("/glossaries/%s", glossaryID), nil, glossary); err != nil {
		return nil, err
	}
	return glossary, nil
}

func (c *Client) DeleteGlossary(ctx context.Context, glossaryID string) error {
	if _, err := c.Delete(ctx, fmt.Sprintf("/glossaries/%s", glossaryID)); err != nil {
		return err
	}
	return nil
}

func (c *Client) GetGlossaryEntries(ctx context.Context, glossaryID string) ([]GlossaryEntry, error) {
	req, err := c.NewRequest(ctx, http.MethodGet, fmt.Sprintf("/glossaries/%s/entries", glossaryID), nil, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "text/tab-separated-values")

	resp, err := c.httpClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	ok := http.StatusOK <= resp.StatusCode && resp.StatusCode < http.StatusMultipleChoices
	if !ok {
		return nil, HandleHTTPError(resp)
	}

	return DecodeGlossaryEntries(resp.Body)
}

type GlossaryEntry struct {
	Source string
	Target string
}

type GlossaryEntries []GlossaryEntry

func NewGlossaryEntries(args ...GlossaryEntry) GlossaryEntries {
	entries := make(GlossaryEntries, 0, len(args))
	for _, arg := range args {
		entries = append(entries, arg)
	}
	return entries
}

func (ges GlossaryEntries) EncodeGlossaryEntries(format EntriesFormat) string {
	out := make([]string, 0, len(ges))
	for _, entry := range ges {
		out = append(out, fmt.Sprintf("%s%s%s", entry.Source, format.Delimiter(), entry.Target))
	}
	return strings.Join(out, "\n")
}

func DecodeGlossaryEntries(r io.Reader) (GlossaryEntries, error) {
	reader := csv.NewReader(r)
	reader.Comma = '\t'
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	entries := make(GlossaryEntries, 0, len(records))
	for _, record := range records {
		entries = append(entries, GlossaryEntry{Source: record[0], Target: record[1]})
	}

	return entries, nil
}

type EntriesFormat interface {
	String() string
	Delimiter() string
}

type (
	csvFormat struct{}
	tsvFormat struct{}
)

var (
	CSV csvFormat
	TSV tsvFormat
)

func (cf csvFormat) String() string {
	return "csv"
}

func (cf csvFormat) Delimiter() string {
	return ","
}

func (tf tsvFormat) String() string {
	return "tsv"
}

func (tf tsvFormat) Delimiter() string {
	return "\t"
}
