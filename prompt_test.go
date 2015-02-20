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
		It("Should default to using the default line editor", func() {
			Expect(prompt.lineEditor).To(BeAssignableToTypeOf(DefaultLineEditor{}))
		})

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
			lineEditor.err = errors.New("Close error")
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
			lineEditor.response = "response"
			lineEditor.err = errors.New("OK")

			prompt.SetLineEditor(lineEditor)
			response, err := prompt.NextResponse()
			Expect(err).To(MatchError("OK"))
			Expect(response).To(Equal("response"))
		})
	})
})
