# ollama-tui

`ollama-tui` is a lightweight terminal user interface built for internal company workstations to interact with local Ollama instances. It enables staff to run prompts, translate internal memos, and summarize attachments locally without uploading sensitive company data to external cloud AI providers.

![TUI Interface Preview](docs/screenshot-placeholder.png)

## Features

- **Interactive TUI Chat**: Real-time response streaming for direct conversations.
- **Local Model Selector**: Automatically detects and lists all models installed on the local Ollama instance.
- **Internal Presets**: Quick access to pre-configured company prompts (memo translation, standard document summarization).
- **File Attachments**: Pipe text files directly into the utility via stdin or attach files using flags for quick contextual analysis.

## Prerequisites

- Go 1.26 or higher
- Local [Ollama](https://ollama.com) daemon running on default port (`http://localhost:11434`)

## Installation & Setup

Clone the repository and build the binary using the Go 1.26 toolchain:

```bash
git clone https://github.com/internal-corp/ollama-tui.git
cd ollama-tui
go build -o ollama-tui .
```

Move the compiled binary to your system PATH or execute it directly:

```bash
./ollama-tui
```

## Usage

### Interactive Mode
Run without arguments to launch the standard interface:
```bash
./ollama-tui
```

### Document Summarization via Pipe
Pipe internal text documents directly for instant analysis:
```bash
cat surat_jalan_2026.txt | ./ollama-tui --preset=summarize
```

### File Flag Attachment
Attach specific documents using the `-f` flag:
```bash
./ollama-tui -f ./draft_nota_dinas.txt --model=llama3
```

## Configuration

By default, `ollama-tui` connects to `http://localhost:11434`. You can override the host by setting the standard Ollama host environment variable:

```bash
export OLLAMA_HOST=http://127.0.0.1:11434
```
