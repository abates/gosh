package gosh

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type testEditor struct {
	prompts   chan string
	responses chan string
}

func (t testEditor) Prompt(p string) (string, error) {
	t.prompts <- p
	return <-t.responses, nil
}

func (t testEditor) sendResponse(response string) {
	t.responses <- response
}

func newTestEditor() *testEditor {
	return &testEditor{
		make(chan string, 10),
		make(chan string, 10),
	}
}

func (t testEditor) GetPrompt() string {
	return "prompt> "
}

var _ = Describe("Shell", func() {
	var editor *testEditor
	var shell *Shell

	BeforeEach(func() {
		editor = newTestEditor()

		shell = NewShell(CommandMap{})
		shell.SetLineEditor(editor)
		shell.SetPrompter(editor)
		shell.Exec()
	})

	Describe("The shell prompt", func() {
		It("should display the prompt built from the prompter interface", func() {
			Expect(<-editor.prompts).To(Equal("prompt> "))
		})
	})

	Describe("Executing a multi-level command", func() {
		var command *complexCommand
		BeforeEach(func() {
			command = newComplexCommand()
			shell.AddCommand("cmd", command)
		})

		It("Should call the top level command if no sub commands are provided", func() {
			editor.sendResponse("cmd\n")
			<-editor.prompts
			<-editor.prompts
			Expect(command.executed).To(BeTrue())
		})

		It("Should call the next level command if a valid next-level command is provided", func() {
			editor.sendResponse("cmd subCmd1\n")
			<-editor.prompts
			<-editor.prompts
			Expect(command.executed).To(Equal(false))
			Expect(command.subCommands["subCmd1"].executed).To(BeTrue())
			Expect(command.subCommands["subCmd1"].arguments).To(Equal([]Argument{}))
		})

		It("Should provide an additional argument to the next level command when given", func() {
			editor.sendResponse("cmd subCmd1 arg1\n")
			<-editor.prompts
			<-editor.prompts
			Expect(command.executed).To(Equal(false))
			Expect(command.subCommands["subCmd1"].executed).To(BeTrue())
			Expect(command.subCommands["subCmd1"].arguments).To(Equal([]Argument{Argument("arg1")}))
		})

		It("Should provide all additional arguments to the next level command when given", func() {
			editor.sendResponse("cmd subCmd1 arg1 arg2\n")
			<-editor.prompts
			<-editor.prompts
			Expect(command.subCommands["subCmd1"].arguments).To(Equal([]Argument{
				Argument("arg1"),
				Argument("arg2"),
			}))
		})
	})
})
