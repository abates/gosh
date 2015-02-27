/**
 * Copyright 2015 Andrew Bates
 *
 * Licensed under the Apache License, Version 2.0
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

// Shell is the foundation for Gosh
//
// A Shell provides a way to prompt users for command input and then execute
// those commands.  It includes line editing, history and command completion.
type Shell struct {
	prompt      Prompt
	commands    CommandMap
	errorWriter io.Writer
}

// SetPrompter overrides the prompter when the DefaultPrompt is being used
//
// A nil prompt will generate the ErrNilPrompter error and trying to override
// the prompter when the DefaultPrompt is not being used generates
// ErrDefaultPrompter
func (shell *Shell) SetPrompter(prompter Prompter) error {
	if prompter == nil {
		return ErrNilPrompter
	}

	if defaultPrompter, ok := shell.prompt.(*DefaultPrompt); ok {
		defaultPrompter.SetPrompter(prompter)
		return nil
	}
	return ErrDefaultPrompter
}

// SetPrompt overrides the Shell' Prompt implementation
//
// Shell is intialized with a DefaultPrompt to prompt the user and gather
// responses.  This includes command completion and history.  However, the
// Prompt can be overridden by a different implementation.  If a nil Prompt is
// given, then ErrNilPrompt is returned
func (shell *Shell) SetPrompt(prompt Prompt) error {
	if prompt == nil {
		return ErrNilPrompt
	}
	shell.prompt = prompt
	return nil
}

// SetErrorWriter overrides the error stream
//
// Shell defaults to use os.Stderr for error messages.  This can be overridden
// with a non-nil io.Writer.  A nil writer generates the ErrNilWriter error
func (shell *Shell) SetErrorWriter(writer io.Writer) error {
	if writer == nil {
		return ErrNilWriter
	}
	shell.errorWriter = writer
	return nil
}

// NewShell returns a fully initialized Shell for the given CommandMap
func NewShell(commands CommandMap) *Shell {
	return &Shell{
		prompt:      NewDefaultPrompt(commands),
		commands:    commands,
		errorWriter: os.Stderr,
	}
}

// Exec starts the Shell prompt/execute loop.
//
// Exec returns upon io.EOF in the input stream
func (shell *Shell) Exec() {
	if prompt, ok := shell.prompt.(Closeable); ok {
		defer prompt.Close()
	}

	for {
		input, err := shell.prompt.NextResponse()

		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Fprintf(shell.errorWriter, "%v\n", err)
			continue
		}

		fields := strings.Fields(input)
		if len(fields) > 0 {
			err = shell.commands.Exec(fields)

			if err != nil {
				fmt.Fprintf(shell.errorWriter, "%v\n", err)
			}
		}
	}
}
