package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
)

// Preset represents approved compliance templates matching internal requirements
type Preset struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	SystemPrompt string `json:"system_prompt"`
}

// Config holds the application level configuration including Ollama settings and Presets.
