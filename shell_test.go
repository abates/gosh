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
	"bufio"
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io"
)

type testEditor struct {
	prompts   chan string
	responses chan string
	err       error
}

func (t testEditor) Prompt(p string) (string, error) {
	t.prompts <- p
	if t.err != nil {
		return "", t.err
	} else {
		return <-t.responses, nil
	}
}

func (t testEditor) sendResponse(response string) {
	t.responses <- response
}

func newTestEditor() *testEditor {
	return &testEditor{
		make(chan string, 10),
		make(chan string, 10),
		nil,
	}
}

func (t testEditor) GetPrompt() string {
	return "prompt> "
}

type errorCommand struct{}

func (e errorCommand) SubCommands() CommandMap {
	return nil
}

func (e errorCommand) Exec(arguments []Argument) error {
	return errors.New("This command failed to execute")
}

var _ = Describe("Shell", func() {
	Describe("interaction", func() {
		var editor *testEditor
		var shell *Shell
		var stderr, stdout *bufio.Reader

		BeforeEach(func() {
			stderr_r, stderr_wr := io.Pipe()
			stdout_r, stdout_wr := io.Pipe()
			stderr = bufio.NewReader(stderr_r)
			stdout = bufio.NewReader(stdout_r)

			editor = newTestEditor()

			shell = NewShell(CommandMap{
				"error": errorCommand{},
			})

			shell.SetLineEditor(editor)
			shell.SetStderr(stderr_wr)
			shell.SetStdout(stdout_wr)
			shell.Exec()
		})

		Describe("The Default prompt", func() {
			It("should display the default prompt", func() {
				Expect(<-editor.prompts).To(Equal("> "))
			})

			It("Should display the prompt from a customized prompter", func() {
				shell.SetPrompter(editor)
				Expect(<-editor.prompts).To(Equal("prompt> "))
			})

			It("should print an error if the prompt encounters an error", func() {
				errorString := "You can't eat, I have no cheeseburgers!"
				editor.err = errors.New(errorString)
				editor.responses <- "eat\n"
				line, _, _ := stderr.ReadLine()
				Expect(string(line)).To(Equal(errorString))
			})

			It("should print an error if the command returns an error", func() {
				editor.responses <- "error\n"
				line, _, _ := stderr.ReadLine()
				Expect(string(line)).To(Equal("This command failed to execute"))
			})
		})
	})

	Describe("parsing the line into arguments", func() {
		It("Should parse a line into fields and convert them into an array of Argument", func() {
			args := getArguments("cmd arg1 arg2 arg3")
			Expect(args).To(Equal([]Argument{
				Argument("cmd"),
				Argument("arg1"),
				Argument("arg2"),
				Argument("arg3"),
			}))
		})
	})
})
