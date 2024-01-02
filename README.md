# askgpt

`askgpt` is a GPT-based chatbot command line tool that sends a prompt to multiple LLMs (e.g. OpenAI's `gpt-4` and
Bard's `text-bison-001`) and displays all the responses. Conversational memory allows follow-up questions to previous
questions by sending context as part of the prompt.

## Example session
This demonstrates the ability to ask follow-up questions using session history:
```bash
$ askgpt --gpt4 "who won the superbowl in 2006?"
Q: who won the superbowl in 2006?
A(gpt-4-1106-preview): The Pittsburgh Steelers won Super Bowl XL (40) in 2006, defeating the Seattle Seahawks with a score of 21-10. The game was played on February 5, 2006, at Ford Field in Detroit, Michigan. This victory earned the Steelers their fifth Super Bowl title at the time.

$ askgpt --gpt4 "and 2008?"
Q: and 2008?
A(gpt-4-1106-preview): In 2008, the New York Giants won Super Bowl XLII (42). They defeated the New England Patriots with a score of 17-14. The game took place on February 3, 2008, at the University of Phoenix Stadium in Glendale, Arizona. This win was particularly notable as it ended the Patriots' quest for a perfect season after they had won all their regular season and playoff games leading up to the Super Bowl.

$ askgpt "and 2009?"
Q: and 2009?
A(bard): Pittsburgh Steelers
A(gpt-4-1106-preview): In 2009, the Pittsburgh Steelers won Super Bowl XLIII (43). They defeated the Arizona Cardinals with a score of 27-23. The game was played on February 1, 2009, at Raymond James Stadium in Tampa, Florida. With this victory, the Steelers secured their sixth Super Bowl title, which at the time made them the franchise with the most Super Bowl wins.
A(claude-2.1):  The Pittsburgh Steelers won Super Bowl XLIII (43) in 2009, defeating the Arizona Cardinals by a score of 27-23. The game was played on February 1, 2009 at Raymond James Stadium in Tampa, Florida. This was the Steelers' sixth Super Bowl championship in franchise history.

So to summarize:

2006 Super Bowl: Pittsburgh Steelers
2008 Super Bowl: New York Giants
2009 Super Bowl: Pittsburgh Steelers
```


## Usage
```
Usage: askgpt [OPTIONS] PROMPT
    OPTIONS:
        --info          Show info and above logs.
        --skip-history  Skip writing to the history.
        --gpt4          Use GPT-4 model.
        --bard          Use Bard model.
        --claude        Use Claude model.
        --gpt4          Shortcut to set LLM_MODELS=gpt-4
        --bard          Shortcut to set LLM_MODELS=bard
        --info          Set log verbosity to info level
    PROMPT              A string prompt to send to the GPT models.

Environment variables:
  PROMPT_PREFIX       A prefix to add to the prompt read from STDIN.
  OPENAI_API_KEY      API key for OpenAI
  ANTHROPIC_API_KEY   API key for Anthropic
  BARDAI_API_KEY      API key for Bard AI
  LLM_MODELS          Comma-separated list of LLM models
  MAX_TOKENS          Maximum number of tokens for a prompt
  HISTORY_LINE_COUNT  How many lines of history context should be considered to add to the prompt

Examples:
  askgpt "Generate go code to iterate over a list"
  askgpt "Refactor that generated code to be in a function named Scan()"
  cat main.go | PROMPT_PREFIX="Generate a code review: " askgpt
  askgpt --skip-history "Generate go code to iterate over a list"
  askgpt --gpt4 "What is the meaning of life?"
  askgpt --bard "Tell me a story about a robot."
  askgpt --claude "Explain quantum computing in simple terms."
```

The program uses a history to provide prior question context to the current prompt, enabling a short-term memory across
sessions that feels like a conversational chat.

## Available Models

`askgpt` supports the following LLM models:

- `openai.GPT3Dot5Turbo`
- `openai.GPT4`
- `openai.GPT4TurboPreview`
- `google.bard`
- `anthropic.claude21`


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
