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
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type testPrompt struct {
	responses chan string
	err       error
}

func (t testPrompt) NextResponse() (string, error) {
	if t.err != nil {
		return "", t.err
	} else {
		return <-t.responses, nil
	}
}

func (t testPrompt) sendResponse(response string) {
	t.responses <- response
}

func newTestPrompt() *testPrompt {
	return &testPrompt{
		make(chan string, 10),
		nil,
	}
}

type errorCommand struct{}

func (e errorCommand) Completions() []string {
	return nil
}

func (e errorCommand) Exec(arguments []string) error {
	return errors.New("This command failed to execute")
}

var _ = Describe("Shell", func() {
	var shell *Shell
	BeforeEach(func() {
		shell = NewShell(CommandMap{})
	})

	Describe("the prompt", func() {
		It("Should use DefaultPrompt", func() {
			Expect(shell.prompt).To(BeAssignableToTypeOf(&DefaultPrompt{}))
		})

		It("Should allow overriding the default prompt", func() {
			tp := newTestPrompt()
			err := shell.SetPrompt(tp)
			Expect(err).To(BeNil())
			Expect(shell.prompt).To(Equal(tp))
		})

		It("Should prohibit setting a nil prompt", func() {
			oldPrompt := shell.prompt
			err := shell.SetPrompt(nil)
			Expect(err).To(MatchError(ErrNilPrompt))
			Expect(shell.prompt).To(Equal(oldPrompt))
		})
	})

	Describe("interaction", func() {
		var prompt *testPrompt
		var shell *Shell

		BeforeEach(func() {
			prompt = newTestPrompt()

			shell = NewShell(CommandMap{
				"error": errorCommand{},
			})

			shell.SetPrompt(prompt)
			go shell.Exec()
		})

		/*Describe("The Default prompt", func() {
			It("should display the default prompt", func() {
				Expect(<-prompt.prompts).To(Equal("> "))
			})

			It("Should display the prompt from a customized prompter", func() {
				shell.Prompt.SetPrompter(func() {
					return "prompt> "
				})
				Expect(<-prompt.prompts).To(Equal("prompt> "))
			})

			It("should print an error if the prompt encounters an error", func() {
				errorString := "You can't eat, I have no cheeseburgers!"
				prompt.err = errors.New(errorString)
				prompt.responses <- "eat\n"
				line, _, _ := stderr.ReadLine()
				Expect(string(line)).To(Equal(errorString))
			})

			It("should print an error if the command returns an error", func() {
				prompt.responses <- "error\n"
				line, _, _ := stderr.ReadLine()
				Expect(string(line)).To(Equal("This command failed to execute"))
			})
		})*/
	})
})
