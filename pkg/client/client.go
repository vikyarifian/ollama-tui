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
