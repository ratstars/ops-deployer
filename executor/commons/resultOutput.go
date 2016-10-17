package commons

import ()

const (
	OUTPUT_INFO  = "INFO"
	OUTPUT_ERROR = "ERROR"
)

type ResultOutput interface {
	Type() string
	String() string
}

type StderrOutput struct {
	content string
}

func (out StderrOutput) Type() string {
	return OUTPUT_ERROR
}

func (out StderrOutput) String() string {
	return out.content
}

func NewStderrOutput(content string) ResultOutput {
	return StderrOutput{
		content: content,
	}
}

type StdoutOutput struct {
	content string
}

func (out StdoutOutput) Type() string {
	return OUTPUT_INFO
}

func (out StdoutOutput) String() string {
	return out.content
}

func NewStdoutOutput(content string) ResultOutput {
	return StdoutOutput{
		content: content,
	}
}

type CustomOutput struct {
	types   string
	content string
}

func (out CustomOutput) Type() string {
	return out.types
}

func (out CustomOutput) String() string {
	return out.content
}

func NewCustomOutput(types, content string) ResultOutput {
	return CustomOutput{
		types:   types,
		content: content,
	}
}

func Merge(a, b ResultOutput) ResultOutput {
	content := a.String() + b.String()
	types := a.Type()
	return NewCustomOutput(types, content)
}
