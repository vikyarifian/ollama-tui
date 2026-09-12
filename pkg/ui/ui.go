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
