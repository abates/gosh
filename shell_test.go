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
	"os"
)

type errorCommand struct{}

func (e errorCommand) Completions() []string {
	return nil
}

func (e errorCommand) Exec(arguments []string) error {
	return errors.New("This command failed to execute")
}

var _ = Describe("Shell", func() {
	var shell *Shell
	var commands CommandMap
	BeforeEach(func() {
		commands = CommandMap{}
		shell = NewShell(commands)
	})

	Describe("the prompt", func() {
		It("Should use DefaultPrompt", func() {
			Expect(shell.prompt).To(BeAssignableToTypeOf(&DefaultPrompt{}))
		})

		It("Should allow overriding the default prompt", func() {
			tp := newTestPrompt()
			Expect(shell.SetPrompt(tp)).To(Succeed())
			Expect(shell.prompt).To(Equal(tp))
		})

		It("Should prohibit setting a nil prompt", func() {
			oldPrompt := shell.prompt
			err := shell.SetPrompt(nil)
			Expect(err).To(MatchError(ErrNilPrompt))
			Expect(shell.prompt).To(Equal(oldPrompt))
		})

		It("Should prohibit setting a nil prompter", func() {
			err := shell.SetPrompter(nil)
			Expect(err).To(MatchError(ErrNilPrompter))
		})
	})

	Describe("errors", func() {
		It("Should use os.Stderr by default", func() {
			Expect(shell.errorWriter).To(Equal(os.Stderr))
		})

		It("Should allow overriding the error writer", func() {
			_, pwr := io.Pipe()
			Expect(shell.SetErrorWriter(pwr)).To(Succeed())
			Expect(shell.errorWriter).To(Equal(pwr))
		})

		It("Should prohibit setting the error writer to nil", func() {
			Expect(shell.SetErrorWriter(nil)).To(MatchError(ErrNilWriter))
			Expect(shell.errorWriter).To(Equal(os.Stderr))
		})
	})

	Describe("Exec", func() {
		var prompt *testPrompt
		var stderr *bufio.Reader

		BeforeEach(func() {
			prompt = newTestPrompt()
			shell.SetPrompt(prompt)
			stderrR, stderrWr, _ := os.Pipe()
			stderr = bufio.NewReader(stderrR)
			shell.SetErrorWriter(stderrWr)
		})

		It("Should always call Close() on closeable prompts when exiting", func() {
			prompt.lineEditor.response <- testResponse{
				"",
				io.EOF,
			}
			shell.Exec()
			Expect(prompt.closeCalled).To(BeTrue())
		})

		It("Should display an error if the prompt returns an error", func() {
			prompt.lineEditor.addResponse("", errors.New("Prompt Error"))
			prompt.lineEditor.addResponse("", io.EOF)
			shell.Exec()
			line, _, _ := stderr.ReadLine()
			Expect(string(line)).To(Equal("Prompt Error"))
		})

		Describe("executing a command", func() {
			var command *testCommand

			BeforeEach(func() {
				command = newTestCommand()
				commands.Add("test", command)
			})

			It("Should execute a command with no arguments", func() {
				prompt.lineEditor.addResponse("test", nil)
				prompt.lineEditor.addResponse("", io.EOF)
				shell.Exec()
				Expect(command.executed).To(BeTrue())
				Expect(command.arguments).To(Equal([]string{"test"}))
			})

			It("Should execute a command with some arguments", func() {
				prompt.lineEditor.addResponse("test arg1 arg2", nil)
				prompt.lineEditor.end()
				shell.Exec()
				Expect(command.executed).To(BeTrue())
				Expect(command.arguments).To(Equal([]string{"test", "arg1", "arg2"}))
			})

			It("Should display an error if the command execution fails", func() {
				command.execErr = errors.New("command error")
				prompt.lineEditor.addResponse("test", nil)
				prompt.lineEditor.end()
				shell.Exec()
				line, _, _ := stderr.ReadLine()
				Expect(string(line)).To(Equal("command error"))
			})
		})
	})
})
