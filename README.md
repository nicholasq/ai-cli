# AI CLI

A simple Golang CLI program that wraps Ollama, enhancing responses with optional 'chain of thought' prompting.

![Version](https://img.shields.io/badge/version-0.1-blue)
![License](https://img.shields.io/badge/license-MIT-green)

## Overview

AI CLI is a command-line interface tool that interacts with AI models through Ollama. Its main feature is the optional implementation of 'chain of thought' prompting, which enhances the quality and depth of AI responses.

## Features

- Interact with Ollama-based AI models
- Optional 'chain of thought' reasoning
- [Not yet] Web search and function calling capabilities
- Debug mode for detailed logging

## Installation

### Prerequisites

- Go 1.23 or later
- Ollama installed and running

### Steps

1. Clone the repository:
   ```
   git clone https://github.com/nicholasq/ai-cli.git
   ```
2. Navigate to the project directory:
   ```
   cd ai-cli
   ```
3. Build the project:
   ```
   go build
   ```

## Usage

Run queries using the `run` command:

```
ai run "What is the capital of France?"
```

### Options

- `--cot`: Enable chain of thought reasoning
- `--debug`: Print debugging logs
- `--model`: Allow choice of LLM

### Examples

```
ai run "Explain the process of photosynthesis"
ai run --model mistral "Explain the process of photosynthesis"
ai run --cot "Explain the process of photosynthesis" # This should lead to a more accurate answer.
ai run --debug "What are the main causes of climate change?"
```

## How It Works

The 'chain of thought' feature works by breaking down the query process into multiple stages:

1. Initial thought generation
2. Reflection and improvement on the initial thoughts
3. Final response generation

This process allows for more nuanced and comprehensive responses.

## Project Structure

- `cmd/`: Command definitions using Cobra
- `internal/`: Internal packages
  - `ai/`: AI client implementations
  - `config/`: Configuration structures

## Dependencies

- [Cobra](https://github.com/spf13/cobra): CLI framework
- [Ollama](https://ollama.ai/): AI model integration

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Troubleshooting

If you encounter issues:

1. Ensure Ollama is running and accessible
2. Check your Go version is 1.23 or later
3. Verify all dependencies are correctly installed

## Development Status

This project is currently in alpha. Features and APIs may change.

## Future Plans

- Support for websearch
- Support for function calling
- Enhanced error handling and recovery

## Reporting Issues

If you encounter any bugs or have feature requests, please open an issue on the GitHub repository.

## Updating

To update the AI CLI tool, pull the latest changes from the repository and rebuild:

```
git pull
go build
```

For any questions or additional information, please open an issue on the GitHub repository.
