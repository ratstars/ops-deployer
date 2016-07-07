// sshbash包是一个远程SSH SHELL执行器，它会连接远端SSH服务器，并在运行上边的bash
// 并在bash上执行指令
package sshbash

import (
	"github.com/ratstars/ops-deployer/view"
)

// SSH的登录信息
type SSHLoginInfo struct{
	ip string
	username string
	password string
}

type SshExecutor struct {
	//是否已经登录, 执行Init成功之后, isLogin会变成true, 否则将为false
	isLogin bool
	//执行器的登录信息, 包括SSH登录需要的所有信息
	LoginInfo SSHLoginInfo
}

//初始化SSH执行器
func (ssh *SshExecutor) Init() {
	//
}