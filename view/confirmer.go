package view

import (
	"github.com/ratstars/ops-deployer/executor/commons"
	"github.com/ratstars/ops-deployer/script"
)

type Confirmer interface {
	Confirm(info string) bool
	DisplayAndPause(info string)
}

type View interface {
	NotifyDisplay(cmd *script.CommandDescriber, result []commons.ResultOutput, isOK bool)
}
