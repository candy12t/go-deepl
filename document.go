package deepl

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type TranslateDocumentOption struct {
	SourceLang string
	Filename   string
	Formality  string
	GlossaryID string
}

type DocumentInfo struct {
	DocumentID  string `json:"document_id"`
	DocumentKey string `json:"document_key"`
}

type DocumentStatus struct {
	DocumentID       string `json:"document_id"`
	Status           string `json:"status"`
	SecondsRemaining int    `json:"seconds_remaining"`
	BilledCharacters int    `json:"billed_characters"`
	ErrorMessage     string `json:"error_message,omitempty"`
}

func (c *Client) UploadTranslateDocument(ctx context.Context, filePath, targetLang string, option *TranslateDocumentOption) (*DocumentInfo, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var body bytes.Buffer
	w := multipart.NewWriter(&body)

	fileWriter, err := w.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(fileWriter, f); err != nil {
		return nil, err
	}

	if err := w.WriteField("target_lang", targetLang); err != nil {
		return nil, err
	}

	if err := writeFieldOption(w, option); err != nil {
		return nil, err
	}

	w.Close()

	req, err := c.NewRequest(ctx, http.MethodPost, "/document", nil, &body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	documentInfo := new(DocumentInfo)
	if _, err := c.Do(req, documentInfo); err != nil {
		return nil, err
	}
	return documentInfo, nil
}

// TODO: use reflect
func writeFieldOption(w *multipart.Writer, option *TranslateDocumentOption) error {
	if option != nil {
		return nil
	}

	if option.SourceLang != "" {
		if err := w.WriteField("source_lang", option.SourceLang); err != nil {
			return err
		}
	}
	if option.Filename != "" {
		if err := w.WriteField("filename", option.Filename); err != nil {
			return err
		}
	}
	if option.Formality != "" {
		if err := w.WriteField("formality", option.Formality); err != nil {
			return err
		}
	}
	if option.GlossaryID != "" {
		if err := w.WriteField("glossary_id", option.GlossaryID); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) CheckDocumentStatus(ctx context.Context, documentID, documentKey string) (*DocumentStatus, error) {
	params := struct {
		DocumentKey string `json:"document_key"`
	}{
		DocumentKey: documentKey,
	}
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	status := new(DocumentStatus)
	if _, err := c.Post(ctx, fmt.Sprintf("/document/%s", documentID), bytes.NewReader(body), status); err != nil {
		return nil, err
	}
	return status, nil
}

func (c *Client) DownloadTranslatedDocument(ctx context.Context, documentID, documentKey string) ([]byte, error) {
	params := struct {
		DocumentKey string `json:"document_key"`
	}{
		DocumentKey: documentKey,
	}
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := c.NewRequest(ctx, http.MethodPost, fmt.Sprintf("/document/%s/result", documentID), nil, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	ok := http.StatusOK <= resp.StatusCode && resp.StatusCode < http.StatusMultipleChoices
	if !ok {
		return nil, HandleHTTPError(resp)
	}

	return io.ReadAll(resp.Body)
}
