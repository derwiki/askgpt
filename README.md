# askgpt

`askgpt` is a Golang application that fans a ChatGPT prompt out to many LLMs and
presents each response to the user. It reads a prompt from the first command line argument
(e.g. a quoted string) or piped from STDIN. PROMPT_PREFIX can be used to add a
prefix to what is being piped in over STDIN, e.g. 'review this source code: '.
Only return syntactically valid markdown.

## Setup

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

## Available Models

`askgpt` supports the following GPT-3 models:

- `openai.GPT3Dot5Turbo`
- `openai.GPT4`
- `text-davinci-003`

It also supports the `bard` model from Google.

## Output

The program will output each response from the models along with their respective model names.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
