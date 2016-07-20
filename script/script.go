package script

import (

)

type Scripter struct{
	Executors []ExecutorDescriber
	Commands []CommandDescriber
}

type ExecutorDescriber struct {
	Name string
	Type string
	Args string
}

type CommandDescriber struct {
	// 当isCommnet为true时, 只有Command有效的, Command的内容为提示内容
	IsComment bool
	ExecutorName string
	Command string
	Timeout int
	ExpectRegular string
	UnexpectRegular string
}