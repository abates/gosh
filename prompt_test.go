package gosh

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type testLineEditor struct {
	closeError  bool
	closeCalled bool
}

func (t *testLineEditor) Prompt(string) (string, error) { return "", nil }
func (t *testLineEditor) Close() error {
	t.closeCalled = true
	if t.closeError {
		return errors.New("Close error")
	}
	return nil
}

type nonCloseableLineEditor struct{}

func (t *nonCloseableLineEditor) Prompt(string) (string, error) { return "", nil }

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
			err := prompt.SetPrompter(func() string { return "custom prompt> " })
			Expect(prompt.prompter()).To(Equal("custom prompt> "))
			Expect(err).To(BeNil())
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
			err := prompt.SetLineEditor(&testLineEditor{})
			Expect(err).To(BeNil())
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
			lineEditor = &testLineEditor{false, false}
			prompt.SetLineEditor(lineEditor)
		})

		It("Should be able to close a closeable line editor", func() {
			err := prompt.Close()
			Expect(err).To(BeNil())
			Expect(lineEditor.closeCalled).To(BeTrue())
		})

		It("Should return an error when closing the line editor returns an error", func() {
			lineEditor.closeError = true
			err := prompt.Close()
			Expect(err).To(MatchError("Close error"))
		})

		It("Should return nil when closing and the line editor is not closeable", func() {
			prompt.SetLineEditor(&nonCloseableLineEditor{})
			err := prompt.Close()
			Expect(err).To(BeNil())
		})
	})
})
