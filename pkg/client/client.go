package client

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Message represents a single message in the chat history,
// matching the role/content structure required by Ollama API.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest represents the request payload sent to Ollama's /api/chat endpoint.
type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

// ChatResponse represents a single streamed chunk or the final response
// returned from Ollama's /api/chat endpoint.
type ChatResponse struct {
	Model     string    `json:"model"`
	CreatedAt string    `json:"created_at"`
	Message   Message   `json:"message"`
	Done      bool      `json:"done"`
}

// ModelInfo represents the metadata of a pulled model in Ollama.
type ModelInfo struct {
	Name string `json:"name"`
}

// TagsResponse represents the response payload from Ollama's /api/tags endpoint.
type TagsResponse struct {
	Models []ModelInfo `json:"models"`
}

// Client is a lightweight HTTP client for interfacing with local Ollama APIs.
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient creates a new instance of Client with a default HTTP client configuration.
func NewClient(baseURL string) *Client {
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}
	baseURL = strings.TrimSuffix(baseURL, "/")
	return &Client{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{},
	}
}

// ListModels fetches the list of local models that have been downloaded into Ollama.
func (c *Client) ListModels(ctx context.Context) ([]string, error) {
	url := fmt.Sprintf("%s/api/tags", c.BaseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request to Ollama failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ollama returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var tagsResp TagsResponse
	if err := json.NewDecoder(resp.Body).Decode(&tagsResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	models := make([]string, 0, len(tagsResp.Models))
	for _, m := range tagsResp.Models {
		models = append(models, m.Name)
	}
	return models, nil
}

// ChatStream initiates a streaming request to /api/chat and invokes the callback
// onChunk for each JSON chunk received.
func (c *Client) ChatStream(ctx context.Context, req ChatRequest, onChunk func(ChatResponse) error) error {
	req.Stream = true
	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/api/chat", c.BaseURL)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("request to Ollama failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("ollama returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var chunk ChatResponse
		if err := json.Unmarshal(line, &chunk); err != nil {
			return fmt.Errorf("failed to decode streaming chunk: %w", err)
		}

		if err := onChunk(chunk); err != nil {
			return fmt.Errorf("callback error: %w", err)
		}

		if chunk.Done {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("reading stream failed: %w", err)
	}

	return nil
}
