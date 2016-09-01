package view

import (
	"fmt"
	"github.com/ratstars/ops-deployer/executor/commons"
	"github.com/ratstars/ops-deployer/script"
	"strings"
)

type ShellView struct{}

func (*ShellView) Confirm(info string) bool {
	var input string
	fmt.Println(info)
	fmt.Println("Please input \"yes\" to confirm running:")
	fmt.Scanln(&input)
	result := strings.TrimSpace(input)
	if "yes" == result {
		return true
	}
	return false
}

func (s *ShellView) DisplayAndPause(info string) {
	var input string
	fmt.Println(info)
	fmt.Print("Please enter to continue...")
	fmt.Scanln(&input)
}

func (*ShellView) NotifyDisplay(cmd *script.CommandDescriber, result []commons.ResultOutput, isOK bool) {
	fmt.Println("===========================")
	fmt.Println("[CMD]", cmd.Command)
	if cmd.ExpectRegular != "" {
		fmt.Println("[EXPECT]", cmd.ExpectRegular)
	}
	if cmd.UnexpectRegular != "" {
		fmt.Println("[UNEXPECT]", cmd.UnexpectRegular)
	}
	if cmd.Timeout > 0 {
		fmt.Println("[IN SECOND]", cmd.Timeout)
	}
	if result != nil {
		for _, line := range result {
			fmt.Printf("[%s]", line.Type())
			fmt.Println(line.String())
		}
	}
	if true == isOK {
		fmt.Println("==========SUCCESS==========")
	} else {
		fmt.Println("==========FAILED===========")
	}
}

func (*ShellView) DisplayInfo(info string) {
	fmt.Println(info)
}
