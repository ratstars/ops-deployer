package controller

import (
	
)

type DefaultController {
	//系统确认器, 当执行器需要进行交互确认时, 会调用这个对象的方法
	Confirmer view.Confirmer
}

//空确认器, 当不给确认器赋值而执行SshExecutor.Init(), 用这个确认器
//这个确认器会在确认时返回false
type invalidConfirmer struct {}

func (*invalidConfirmer) Confirm(info string) bool{
	log.Print("ERROR: Invalid Confirmer be invoke. False be returned.")
	return false
}

func (*invalidConfirmer) DisplayAndPause(info string){
	log.Print("ERROR: Invalid Confirmer be invoke.")
}
