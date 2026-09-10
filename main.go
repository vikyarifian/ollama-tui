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
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func main() {
	ctx := context.Background()

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Warning: Gagal memuat konfigurasi, menggunakan default. Error: %v\n", err)
		cfg = config.GetDefaultConfig()
	}

	// Parse command-line flags
	modelFlag := flag.String("model", "", "Nama model Ollama yang akan digunakan")
	presetFlag := flag.String("preset", "", "Nama atau nomor preset kepatuhan yang akan digunakan")
	inputFlag := flag.String("input", "", "Path file input untuk langsung diproses")
	flag.Parse()

	// Check for piped input (stdin)
	var initialInput string
	stat, err := os.Stdin.Stat()
	if err == nil && (stat.Mode()&os.ModeCharDevice) == 0 {
		data, err := io.ReadAll(os.Stdin)
		if err == nil {
			initialInput = strings.TrimSpace(string(data))
		}
	}

	// If input file flag is provided, read it
	if *inputFlag != "" {
		fileData, err := os.ReadFile(*inputFlag)
		if err != nil {
			fmt.Printf("Error: Gagal membaca file input %s: %v\n", *inputFlag, err)
			os.Exit(1)
		}
		if initialInput != "" {
			initialInput = initialInput + "\n\n" + strings.TrimSpace(string(fileData))
		} else {
			initialInput = strings.TrimSpace(string(fileData))
		}
	}

	// Initialize UI
	u := ui.NewUI()

	// Initialize Ollama Client
	cli := client.NewClient(cfg.OllamaURL)

	// Get available models from Ollama to let user choose
	var models []string
	models, err = cli.ListModels(ctx)
	if err != nil {
		fmt.Printf("Warning: Tidak dapat terhubung ke Ollama di %s: %v\n", cfg.OllamaURL, err)
		fmt.Println("Pastikan Ollama sudah berjalan (ollama serve).")
		// Use default model from config as fallback
		models = []string{cfg.DefaultModel}
	}

	// Determine selected model
	selectedModel := *modelFlag
	if selectedModel == "" {
		selectedModel, err = u.SelectModel(models, cfg.DefaultModel)
		if err != nil {
			fmt.Printf("Error memilih model: %v. Menggunakan default: %s\n", err, cfg.DefaultModel)
			selectedModel = cfg.DefaultModel
		}
	}

	// Determine selected preset
	var selectedPreset config.Preset
	var presetActive bool

	if *presetFlag != "" {
		presetActive = false
		for i, p := range cfg.Presets {
			idxStr := fmt.Sprintf("%d", i+1)
			if *presetFlag == idxStr || strings.Contains(strings.ToLower(p.Name), strings.ToLower(*presetFlag)) {
				selectedPreset = p
				presetActive = true
				break
			}
		}
		if !presetActive {
			fmt.Printf("Warning: Preset '%s' tidak ditemukan. Berjalan tanpa preset.\n", *presetFlag)
		}
	} else {
		var err error
		selectedPreset, presetActive, err = u.SelectPreset(cfg.Presets)
		if err != nil {
			fmt.Printf("Error memilih preset: %v. Berjalan tanpa preset.\n", err)
			presetActive = false
		}
	}

	var finalPreset config.Preset
	if presetActive {
		finalPreset = selectedPreset
	}

	// Start the chat session
	err = u.RunChat(ctx, cli, selectedModel, finalPreset, initialInput)
	if err != nil && err != io.EOF {
		fmt.Printf("Error dalam sesi chat: %v\n", err)
		os.Exit(1)
	}
}
