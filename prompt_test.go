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
	responses   chan string
	lineEditor  *testLineEditor
	closeErr    error
	closeCalled bool
}

func (t *testPrompt) NextResponse() (string, error) {
	return t.lineEditor.Prompt("")
}

func (t *testPrompt) Close() error {
	t.closeCalled = true
	return t.closeErr
}

func newTestPrompt() *testPrompt {
	return &testPrompt{
		make(chan string, 10),
		newTestLineEditor(),
		nil,
		false,
	}
}

var _ = Describe("DefaultPrompt", func() {
	var prompt *DefaultPrompt
	BeforeEach(func() {
		prompt = NewDefaultPrompt(CommandMap{})
	})

	Describe("the prompter", func() {
		It("Should have a default prompter", func() {
			Expect(prompt.prompter()).To(Equal("> "))
		})

		It("Should allow overriding the default prompter", func() {
			Expect(prompt.SetPrompter(func() string { return "custom prompt> " })).To(Succeed())
			Expect(prompt.prompter()).To(Equal("custom prompt> "))
		})

		It("Should not allow setting a nil prompter", func() {
			err := prompt.SetPrompter(nil)
			Expect(err).To(MatchError(ErrNilCallback))
		})

	})

	Describe("the line editor", func() {
		It("Should allow overriding the default line editor", func() {
			Expect(prompt.SetLineEditor(newTestLineEditor())).To(Succeed())
		})

		It("Should not allow setting the line editor to nil", func() {
			oldEditor := prompt.lineEditor
			err := prompt.SetLineEditor(nil)
			Expect(err).To(MatchError(ErrNilLineEditor))
			Expect(prompt.lineEditor).To(Equal(oldEditor))
		})
	})

	Describe("Closing the prompt", func() {
		var lineEditor *testLineEditor
		BeforeEach(func() {
			lineEditor = newTestLineEditor()
			prompt.SetLineEditor(lineEditor)
		})

		It("Should be able to close a closeable line editor", func() {
			Expect(prompt.Close()).To(Succeed())
			Expect(lineEditor.closeCalled).To(BeTrue())
		})

		It("Should return an error when closing the line editor returns an error", func() {
			lineEditor.closeErr = errors.New("Close error")
			err := prompt.Close()
			Expect(err).To(MatchError("Close error"))
		})

		It("Should return nil when closing and the line editor is not closeable", func() {
			prompt.SetLineEditor(&nonCloseableLineEditor{})
			Expect(prompt.Close()).To(Succeed())
		})
	})

	Describe("getting a response", func() {
		It("should return a response captured from LineEditor.Prompt", func() {
			lineEditor := newTestLineEditor()
			lineEditor.addResponse("response", errors.New("OK"))

			prompt.SetLineEditor(lineEditor)
			response, err := prompt.NextResponse()
			Expect(err).To(MatchError("OK"))
			Expect(response).To(Equal("response"))
		})
	})
})
