package view

import (
	"github.com/ratstars/ops-deployer/executor/commons"
)

type Confirmer interface{
	Confirm(info string) bool;
	DisplayAndPause(info string);
}

type View interface{
	NotifyDisplay(result []commons.ResultOutput, isOK bool)
	
}