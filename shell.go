package gosh

import (
	"fmt"
	"os"
	"strings"
)

type Shell struct {
	lineEditor LineEditor
	prompter   Prompter
	commands   CommandMap
	completer  Completer
}

func NewShell(commands CommandMap) *Shell {
	return &Shell{
		nil,
		nil,
		commands,
		*NewCompleter(commands),
	}
}

type LineEditor interface {
	Prompt(string) (string, error)
}

type Prompter interface {
	GetPrompt() string
}

func (shell *Shell) SetLineEditor(lineEditor LineEditor) {
	shell.lineEditor = lineEditor
}

func (shell *Shell) SetPrompter(prompter Prompter) {
	shell.prompter = prompter
}

func (shell *Shell) AddCommand(commandName string, command ShellCommand) error {
	return shell.commands.AddCommand(commandName, command)
}

func (shell *Shell) Exec() {
	go func() {
		for {
			input, err := shell.lineEditor.Prompt(shell.prompter.GetPrompt())
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to read from prompt: %v", err)
				continue
			}

			fields := strings.Fields(input)
			arguments := make([]Argument, len(fields))
			for i, field := range fields {
				arguments[i] = Argument(field)
			}

			err = shell.commands.Exec(arguments)

			if err != nil {
				fmt.Fprintf(os.Stderr, "%v", err)
			}
		}
	}()
}
