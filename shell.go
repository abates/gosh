/**
 * = LICENSE
 *
 * Copyright 2015 Andrew Bates Licensed under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with the
 * License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 * License for the specific language governing permissions and limitations under
 * the License.
 */

package gosh

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type Shell struct {
	lineEditor LineEditor
	prompter   Prompter
	stderr     io.Writer
	stdout     io.Writer
	commands   CommandMap
	completer  Completer
}

func NewShell(commands CommandMap) *Shell {
	return &Shell{
		nil,
		DefaultPrompter{},
		os.Stderr,
		os.Stdout,
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

type DefaultPrompter struct{}

func (p DefaultPrompter) GetPrompt() string {
	return "> "
}

func (shell *Shell) SetLineEditor(lineEditor LineEditor) {
	shell.lineEditor = lineEditor
}

func (shell *Shell) SetPrompter(prompter Prompter) {
	shell.prompter = prompter
}

func (shell *Shell) SetStdout(stdout io.Writer) {
	shell.stdout = stdout
}

func (shell *Shell) SetStderr(stderr io.Writer) {
	shell.stderr = stderr
}

func getArguments(line string) []Argument {
	fields := strings.Fields(line)
	arguments := make([]Argument, len(fields))
	for i, field := range fields {
		arguments[i] = Argument(field)
	}
	return arguments
}

func (shell *Shell) Exec() {
	go func() {
		for {
			input, err := shell.lineEditor.Prompt(shell.prompter.GetPrompt())
			if err != nil {
				fmt.Fprintf(shell.stderr, "%v\n", err)
				continue
			}

			err = shell.commands.Exec(getArguments(input))

			if err != nil {
				fmt.Fprintf(shell.stderr, "%v\n", err)
			}
		}
	}()
}
