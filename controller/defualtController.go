package controller

import (
	"github.com/ratstars/ops-deployer/executor"
	"github.com/ratstars/ops-deployer/executor/commons"
	"github.com/ratstars/ops-deployer/script"
	"github.com/ratstars/ops-deployer/view"
	"log"
	"errors"
)

type DefaultController struct{
	//系统确认器, 当执行器需要进行交互确认时, 会调用这个对象的方法
	Confirmer view.Confirmer
	View view.View
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

func (dc *DefaultController) RunScript(script script.Scripter) error {
	if nil == dc.Confirmer {
		dc.Confirmer = &invalidConfirmer{}
	}
	// 1. 新建各执行器, 并将执行各执行器初始化
	
	// 执行器Map
	executorMap := make(map[string] executor.Executor)
	for _, executorDes := range script.Executors {
		if _, ok := executorMap[executorDes.Name]; false == ok {
			// 没有这个名字的执行器，进行执行器的创建
			executor, err:= executor.CreateExecutor(executorDes.Type, executorDes.Args)
			if err != nil {
				log.Print("Create Executor Failed. ")
				log.Print("Name:", executorDes.Name)
				log.Println("Reason:", err)
				return err
			}
			executorMap[executorDes.Name] = executor
		}
	}
	// 定义defer销毁执行器
	defer func() {
		for _, executor := range executorMap {
			if true == executor.IsReady() {
				executor.Destory()
			}
		}
	}()
	// 初始化各执行器
	for key, executor := range executorMap {
		executor.Init()
		if false == executor.IsReady() {
			log.Println("Executor Init Failed:", key)
			return errors.New("Executor Init Failed.")
		}
	}
	
	// 2. 执行Script
	for _, cmd := range script.Commands {
		if true == cmd.IsComment {
			//是注释行
			dc.Confirmer.DisplayAndPause(cmd.Command)
		} else {
			//非注释
			executor, ok := executorMap[cmd.ExecutorName]
			if false == ok {
				log.Println("Not define executor:", cmd.ExecutorName)
				return errors.New("Executor Not Define")
			}
			if 0 <= cmd.Timeout{
				// 没有设置超时时间, 则将超时时间设置为10s
				cmd.Timeout = 10
			}
			result, err := executor.Execute(cmd.Command, cmd.Timeout)
			if err != nil {
				log.Println("Command Execute Error:",err)
				return err 
			}
			ok = checkResult(result, &cmd)
			// 通知view
			if nil != dc.View {
				dc.View.NotifyDisplay(result, ok)
			}
			if false == ok {
				//执行失败, 不再执行后续指令
				break
			}
		}
	}
	return nil
}

func checkResult(result []commons.ResultOutput, cmd *script.CommandDescriber) bool {
	//TODO
	return false
}