# askgpt

`askgpt` is a GPT-based chatbot command line tool that sends a prompt to multiple LLMs (e.g. OpenAI's `gpt-4` and
Bard's `text-bison-001`) and displays all the responses. Conversational memory allows follow-up questions to previous
questions by sending context as part of the prompt.

## Example session
This demonstrates the ability to ask follow-up questions and query multiple LLMs simultaneously.
```
$ askgpt --gpt4 "who won the superbowl in 2006?"
A(gpt-4):The Super Bowl in 2006 (Super Bowl XL) was won by the Pittsburgh Steelers.

# you can ask follow up questions that recall earlier parts of the conversation
$ askgpt --gpt4 "and in 2008?"
A(gpt-4): The Super Bowl in 2008 (Super Bowl XLII) was won by the New York Giants.

# and retain that conversation across different models
$ askgpt "and in 2009?"
A(text-davinci-003): The Super Bowl in 2009 (Super Bowl XLIII) was won by the Pittsburgh Steelers.
A(gpt-4): The Super Bowl in 2009 (Super Bowl XLIII) was won by the Pittsburgh Steelers.
A(gpt-3.5-turbo): The Super Bowl in 2009 (Super Bowl XLIII) was won by the Pittsburgh Steelers.
A(bard): The Super Bowl in 2009 (Super Bowl XLIII) was won by the Pittsburgh Steelers.
```


## Usage

```
Usage: askgpt [OPTIONS] PROMPT
    OPTIONS:
        --skip-history  Skip reading and writing to the history.
        --gpt4          Shortcut to set LLM_MODELS=gpt-4
        --bard          Shortcut to set LLM_MODELS=bard
        --info          Set log verbosity to info level
    PROMPT              A string prompt to send to the GPT models.

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
  askgpt --skip-history --gpt4 --info "Generate go code to iterate over a list"
```

The program uses a history to provide prior question context to the current prompt, enabling a short-term memory across
sessions that feels like a conversational chat.

## Available Models

`askgpt` supports the following OpenAI GPT models:

- `openai.GPT3Dot5Turbo`
- `openai.GPT4`
- `text-davinci-003`

It also supports the `bard` model from Google.

## Output

The program will output each response from the models along with their respective model names.
The program will also keep a history of all questions asked and answered in `~/.askgpt_history` to provide context going
forward for succcessive questions.

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
