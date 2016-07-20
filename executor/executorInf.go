package executor

import (
	"github.com/ratstars/ops-deployer/executor/commons"
)

// Executor接口，这个接口的实现是各种执行器
type Executor interface {
	// 执行器的初始化, 执行完成之后应用IsReady进行判断是否成功
	Init()
	// 执行器销毁, 销毁后IsReady()方法应返回false值
	Destory()
	// 判断执行器是否已经就绪
	IsReady() bool
	// 执行器执行
	Execute(cmd string, timeout int) ([]commons.ResultOutput, error)
}