package rag

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
)

type ragResp struct {
	Response string `json:"facts"`
}

type ragUpload struct {
	Status string `json:"status"`
	DocID  int64  `json:"doc_id"`
}

type Rag struct {
	addr string
}

func NewRag(addr string) *Rag {
	return &Rag{
		addr: addr,
	}
}

var (
	ErrEmptyPaperIDs          = errors.New("empty paper ids")
	ErrFailedToUploadDocument = errors.New("status code is != 200")
	ErrStatusNotSuccess       = errors.New("returned status isn't succeess")
	ErrMissedID               = errors.New("returned id is missed")
)

const (
	getLLMResponseUrl = "/gen_facts/"
	uploadDocumentUrl = "/upload_document/"
)

func (r *Rag) GetLLMResponse(ctx context.Context, request string, paperIDs []int64) (string, error) {
	if len(paperIDs) == 0 {
		return "", ErrEmptyPaperIDs
	}

	params := url.Values{}
	params.Add("query", request)
	fillParams(params, paperIDs)

	url := r.addr + getLLMResponseUrl + "?" + params.Encode()

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to do get request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var llmResp ragResp
	if err := json.Unmarshal(body, &llmResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal rag response: %w", err)
	}

	return llmResp.Response, nil
}

func fillParams(params url.Values, ids []int64) {
	for _, id := range ids {
		params.Add("doc_ids", strconv.Itoa(int(id)))
	}
}

func (r *Rag) UploadDocument(ctx context.Context, file io.ReadCloser, id int64) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", "data.bin")
	if err != nil {
		return fmt.Errorf("failed to create from file: %w", err)
	}

	if _, err := io.Copy(part, file); err != nil {
		return fmt.Errorf("failed to copy data: %w", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	params := url.Values{}
	params.Add("doc_id", strconv.Itoa(int(id)))
	params.Add("overwrite", "true")

	url := r.addr + uploadDocumentUrl + "?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, "POST", url, body)
	if err != nil {
		return fmt.Errorf("failed to create new request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to do request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ErrFailedToUploadDocument
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read all: %w", err)
	}

	var v ragUpload
	if err := json.Unmarshal(data, &v); err != nil {
		return fmt.Errorf("failed to unmarshal: %w", err)
	}

	if v.Status != "success" {
		return ErrStatusNotSuccess
	}

	if v.DocID != id {
		return ErrMissedID
	}

	return nil
}
