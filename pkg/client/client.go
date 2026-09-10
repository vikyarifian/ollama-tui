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
