# askgpt

`askgpt` is a GPT-based chatbot command line tool that sends a prompt to multiple LLMs (e.g. OpenAI's `gpt-4` and Bard's `text-bison-001`) and displays all the responses. Conversational memory allows follow-up questions to previous questions by sending context as part of the prompt.


## Usage

```
Usage: askgpt [OPTIONS] PROMPT
    OPTIONS:
        --skip-history  Skip reading and writing to the history. This flag can come before or after the PROMPT on the command line.
    PROMPT           A string prompt to send to the GPT models.

Environment variables:
  PROMPT_PREFIX       A prefix to add to the prompt read from STDIN.
  OPENAI_API_KEY      API key for OpenAI
  BARDAI_API_KEY      API key for Bard AI
  LLM_MODELS          Comma-separated list of LLM models
  MAX_TOKENS          Maximum number of tokens for a prompt
  HISTORY_LINE_COUNT  How many lines of history context should be considered to add to the prompt

Examples:
  askgpt "Generate go code to iterate over a list"
  askgpt "Refactor that generated code to be in a function named Scan()"
  cat main.go | PROMPT_PREFIX="Generate a code review: " askgpt
  askgpt --skip-history "Generate go code to iterate over a list"
```

The `--skip-history` flag can be used to skip reading and writing to the history. It can come before or after the `PROMPT` on the command line.

The program uses a history to provide prior question context to the current prompt, enabling a short-term memory across sessions that feels like a conversational chat.

## Available Models

`askgpt` supports the following GPT-3 models:

- `openai.GPT3Dot5Turbo`
- `openai.GPT4`
- `text-davinci-003`

It also supports the `bard` model from Google.

## Output

The program will output each response from the models along with their respective model names.
The program will also keep a history of all questions asked and answered in `~/.askgpt_history` to provide context going forward for succcessive questions.

## Installation

1. Clone the repository:
   ```
   git clone https://github.com/derwiki/askgpt
   ```
2. Build the executable file:
   ```
   go build main.go
   ```
3. Run the program:
   ```
   ./askgpt "your prompt here"
   ```
   Or, you can pipe a prompt to STDIN:
   ```
   echo "your prompt here" | ./askgpt
   ```
   And optionally add a prompt prefix:
   ```
   cat main.go | PROMPT_PREFIX="generate a README.md for this program.\n\nProgram source: ""
   ```
   Send to a subset of models:
   ```
   LLM_MODELS=gpt-4 ./askgpt "how do I read an environment variable?"
   LLM_MODELS=gpt-4,gpt-3.5-turbo ./askgpt "how do I read an environment variable?"
   ```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
