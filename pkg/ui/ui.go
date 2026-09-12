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
