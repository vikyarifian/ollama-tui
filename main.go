package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"ollama-tui/pkg/client"
	"ollama-tui/pkg/config"
	"ollama-tui/pkg/ui"
)

// Preset represents approved compliance templates matching internal requirements
type Preset struct {
	Name         string
	Description  string
	SystemPrompt string
}

// Message represents an individual Ollama chat role and content pair
