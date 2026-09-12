package ui

import ( 
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"ollama-tui/pkg/client"
	"ollama-tui/pkg/config"
)

// ANSI escape codes for coloring terminal output.
const (
	Reset   = "\033[0m"
	Bold    = "\033[1m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35"
	Cyan    = "\033[36m"
	Gray    = "\033[90m"
)

// UI coordinates the user interaction in the terminal.
type UI struct {
	scanner *bufio.Scanner
	out     io.Writer
}

// NewUI creates a new instance of UI using standard input and output.
func NewUI() *UI {
	s := bufio.NewScanner(os.Stdin)

	// Workaround for bufio.Scanner limitation: the default token buffer size limit is 64KB (bufio.MaxScanTokenSize).
	// When operational colleagues pipe large documents, audit logs, or reports for compliance review,
	// the default scanner fails with token too long errors. We increase the buffer capacity to 4MB.
	const maxCapacity = 4 * 1024 * 1024
	buf := make([]byte, 0, 64*1024)
	s.Buffer(buf, maxCapacity)

	return &UI{
		scanner: s,
		out:     os.Stdout,
	}
}

// Prompt displays a prompt message and retrieves a line of input from the user.
func (u *UI) Prompt(label string) (string, error) {
	fmt.Fprintf(u.out, "%s%s%s ", Yellow, label, Reset)
	if !u.scanner.Scan() {
		if err := u.scanner.Err(); err != nil {
			return "", err
		}
		return "", io.EOF
	}
	return strings.TrimSpace(u.scanner.Text()), nil
}

// SelectModel displays the list of local Ollama models and lets the user choose one.
func (u *UI) SelectModel(models []string, current string) (string, error) {
	if len(models) == 0 {
		return current, fmt.Errorf("no models available in Ollama")
	}

	fmt.Fprintln(u.out, "\n"+Bold+"=== PILIH MODEL OLLAMA ==="+Reset)
	for i, model := range models {
		marker := " "
		if model == current {
			marker = "*"
		}
		fmt.Fprintf(u.out, " [%s] %d. %s\n", marker, i+1, model)
	}
	fmt.Fprintf(u.out, " [ ] 0. Gunakan model default (%s)\n", current)

	for {
		input, err := u.Prompt("Masukkan nomor pilihan Anda:")
		if err != nil {
			return current, err
		}

		if input == "" || input == "0" {
			return current, nil
		}

		idx, err := strconv.Atoi(input)
		if err != nil || idx < 1 || idx > len(models) {
			fmt.Fprintln(u.out, Red+"Pilihan tidak valid. Silakan coba lagi."+Reset)
			continue
		}

		return models[idx-1], nil
	}
}

// SelectPreset displays approved compliance and operational presets for selection.
func (u *UI) SelectPreset(presets []config.Preset) (config.Preset, bool, error) {
	fmt.Fprintln(u.out, "\n"+Bold+"=== PILIH PRESET KEPATUHAN / OPERASIONAL ==="+Reset)
	fmt.Fprintln(u.out, Gray+"(Membantu menyusun format sesuai aturan Pak Budi/Compliance)"+Reset)

	for i, p := range presets {
		fmt.Fprintf(u.out, "  %d. %s%s%s\n", i+1, Cyan, p.Name, Reset)
		fmt.Fprintf(u.out, "     %s%s%s\n", Gray, p.Description, Reset)
	}
	fmt.Fprintln(u.out, "  0. Tanpa Preset (Chat Bebas / Polos)")

	for {
		input, err := u.Prompt("Masukkan nomor preset:")
		if err != nil {
			return config.Preset{}, false, err
		}

		if input == "" || input == "0" {
			return config.Preset{}, false, nil
		}

		idx, err := strconv.Atoi(input)
		if err != nil || idx < 1 || idx > len(presets) {
			fmt.Fprintln(u.out, Red+"Pilihan tidak valid. Silakan coba lagi."+Reset)
			continue
		}

		return presets[idx-1], true, nil
	}
}

// RunChat starts the interactive streaming chat loop with the selected Ollama model.
func (u *UI) RunChat(ctx context.Context, cli *client.Client, model string, preset config.Preset, initialInput string) error {
	var history []client.Message

	fmt.Fprintln(u.out, "\n"+Bold+"========================================="+Reset)
	fmt.Fprintf(u.out, " Memulai Sesi Chat (%s%s%s)\n", Green, model, Reset)
	if preset.Name != "" {
		fmt.Fprintf(u.out, " Preset Aktif: %s%s%s\n", Cyan, preset.Name, Reset)
		history = append(history, client.Message{
			Role:    "system",
			Content: preset.SystemPrompt,
		})
	}
	fmt.Fprintln(u.out, Gray+" Ketik 'exit' atau 'quit' untuk keluar, '/clear' untuk mereset chat."+Reset)
	fmt.Fprintln(u.out, Bold+"========================================="+Reset)

	// If there's an initial input (piped or loaded from a file flag)
	if initialInput != "" {
		fmt.Fprintf(u.out, "\n%s[User - File/Piped Input]%s\n%s\n", Green, Reset, initialInput)
		history = append(history, client.Message{
			Role:    "user",
			Content: initialInput,
		})

		reply, err := u.streamResponse(ctx, cli, model, history)
		if err != nil {
			fmt.Fprintf(u.out, "\n%s[Error] Gagal mendapatkan respon: %v%s\n", Red, err, Reset)
		} else {
			history = append(history, client.Message{
				Role:    "assistant",
				Content: reply,
			})
		}
	}

	for {
		userMsg, err := u.Prompt("\nUser >>")
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		if userMsg == "" {
			continue
		}

		lowerMsg := strings.ToLower(userMsg)
		if lowerMsg == "exit" || lowerMsg == "quit" {
			fmt.Fprintln(u.out, Yellow+"Sesi chat diakhiri. Sampai jumpa!"+Reset)
			break
		}

		if lowerMsg == "/clear" {
			history = nil
			if preset.Name != "" {
				history = append(history, client.Message{
					Role:    "system",
					Content: preset.SystemPrompt,
				})
			}
			fmt.Fprintln(u.out, Cyan+"Riwayat chat telah direset!"+Reset)
			continue
		}

		history = append(history, client.Message{
			Role:    "user",
			Content: userMsg,
		})

		reply, err := u.streamResponse(ctx, cli, model, history)
		if err != nil {
			fmt.Fprintf(u.out, "\n%s[Error] Gagal mendapatkan respon: %v%s\n", Red, err, Reset)
		} else {
			history = append(history, client.Message{
				Role:    "assistant",
				Content: reply,
			})
		}
	}

	return nil
}

// streamResponse coordinates streaming slices of responses to the output.
func (u *UI) streamResponse(ctx context.Context, cli *client.Client, model string, history []client.Message) (string, error) {
	fmt.Fprintf(u.out, "\n%sOllama >> %s", Blue, Reset)

	var assistantReply strings.Builder
	req := client.ChatRequest{
		Model:    model,
		Messages: history,
		Stream:   true,
	}

	err := cli.ChatStream(ctx, req, func(chunk client.ChatResponse) error {
		text := chunk.Message.Content
		fmt.Fprint(u.out, text)
		assistantReply.WriteString(text)
		return nil
	})

	fmt.Fprintln(u.out) 

	if err != nil {
		return assistantReply.String(), err
	}

	return assistantReply.String(), nil
}
