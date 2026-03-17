package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"vapi-example-1/internal/model"
)

// VapiClient is an HTTP client for the Vapi REST API.
type VapiClient struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// New creates a new VapiClient.
func New(apiKey, baseURL string) *VapiClient {
	return &VapiClient{
		apiKey:  apiKey,
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// doRequest executes an HTTP request against the Vapi API.
func (c *VapiClient) doRequest(method, path string, body interface{}) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, c.baseURL+path, reqBody)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// CreateAssistant creates a new voice assistant.
func (c *VapiClient) CreateAssistant(req model.CreateAssistantRequest) (*model.Assistant, error) {
	respBody, err := c.doRequest("POST", "/assistant", req)
	if err != nil {
		return nil, err
	}

	var assistant model.Assistant
	if err := json.Unmarshal(respBody, &assistant); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	return &assistant, nil
}

// ListAssistants returns all assistants.
func (c *VapiClient) ListAssistants() ([]model.Assistant, error) {
	respBody, err := c.doRequest("GET", "/assistant", nil)
	if err != nil {
		return nil, err
	}

	var assistants []model.Assistant
	if err := json.Unmarshal(respBody, &assistants); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	return assistants, nil
}

// GetAssistant returns a specific assistant by ID.
func (c *VapiClient) GetAssistant(id string) (*model.Assistant, error) {
	respBody, err := c.doRequest("GET", "/assistant/"+id, nil)
	if err != nil {
		return nil, err
	}

	var assistant model.Assistant
	if err := json.Unmarshal(respBody, &assistant); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	return &assistant, nil
}

// DeleteAssistant deletes an assistant by ID.
func (c *VapiClient) DeleteAssistant(id string) error {
	_, err := c.doRequest("DELETE", "/assistant/"+id, nil)
	return err
}
