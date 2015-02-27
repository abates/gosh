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

type Shell struct {
	prompt      Prompt
	commands    CommandMap
	errorWriter io.Writer
}

func (s *Shell) SetPrompter(prompter Prompter) error {
	if prompter == nil {
		return ErrNilPrompter
	}

	if defaultPrompter, ok := s.prompt.(*DefaultPrompt); ok {
		defaultPrompter.SetPrompter(prompter)
		return nil
	} else {
		return ErrDefaultPrompter
	}
}

func (s *Shell) SetPrompt(prompt Prompt) error {
	if prompt == nil {
		return ErrNilPrompt
	}
	s.prompt = prompt
	return nil
}

func (s *Shell) SetErrorWriter(writer io.Writer) error {
	if writer == nil {
		return ErrNilWriter
	}
	s.errorWriter = writer
	return nil
}

func NewShell(commands CommandMap) *Shell {
	return &Shell{
		prompt:      NewDefaultPrompt(commands),
		commands:    commands,
		errorWriter: os.Stderr,
	}
}

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
