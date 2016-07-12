package sshbash

import (
	"log"
	"golang.org/x/crypto/ssh"
	"bytes"
	"sync"
	"io"
	"time"
	"errors"
)

var promptings = []byte {
	27, 7,27, 7, 27,
}

// SSH的登录信息, Key和string只能有一个生效, Key如果不为nil, 将优先使用
type SSHLoginInfo struct{
	Ip string
	Port string
	Username string
	Password string
	Key []byte
}

// 生成sshClientConfig对象
func (loginInfo *SSHLoginInfo) getClientConfig() (*ssh.ClientConfig, error) {
	if loginInfo.Key != nil {
		// 如果Key不为nil对象
		signer, err := ssh.ParsePrivateKey(loginInfo.Key)
		if err != nil {
			return nil, err
		}
		return &ssh.ClientConfig{
			User: loginInfo.Username,
			Auth: []ssh.AuthMethod {
				ssh.PublicKeys(signer),
			},
		}, nil
	} else {
		return &ssh.ClientConfig{
			User: loginInfo.Username,
			Auth: []ssh.AuthMethod{
		        ssh.Password(loginInfo.Password),
			},
		}, nil
	}
}

type SshExecutor struct {
	//是否已经登录, 执行Init成功之后, isLogin会变成true, 否则将为false
	isLogin bool
	//执行器的登录信息, 包括SSH登录需要的所有信息
	LoginInfo *SSHLoginInfo
	client *ssh.Client
	session *ssh.Session
}

//初始化SSH执行器
func (sshe *SshExecutor) Init() {
	//如果已经初始化, 只是返回
	if sshe.isLogin {
		return
	}
	//检查各参数是否已经OK，如果有成员不完整用默认代替或提示出错
	if nil == sshe.LoginInfo{
		log.Println("ERROR: SSH Login Information is empty. Ssh is not start.")
		return
	}
	
	//登录
	config, err:= sshe.LoginInfo.getClientConfig()
	if err != nil {
		log.Println("ERROR: get Client Config error, errorMsg: ", err)
		return
	}
	var address string
	if sshe.LoginInfo.Port == "" {
		address = sshe.LoginInfo.Ip+":22"
	} else {
		address = sshe.LoginInfo.Ip+":"+sshe.LoginInfo.Port
	}
	sshe.client, err = ssh.Dial("tcp", address, config)
	if err != nil {
		log.Println("ERROR: Failed to dial:", err)
		return
	}
	
	//创建执行SESSION
	sshe.session, err = sshe.client.NewSession()
	if err != nil {
		sshe.clearSessionAndClientWhenError("Failed to create session: ", err)
		return
	}
	prompt_notify := make(chan int)
	defer close(prompt_notify)
	var loginInfoBuffer singleRawWriterBuffer
	sshe.session.Stdout = newDecoratorWriterForNofityer(&loginInfoBuffer, prompt_notify)
	sshe.session.Stderr = &loginInfoBuffer
	stdin, err := sshe.session.StdinPipe()
	if err != nil {
		sshe.clearSessionAndClientWhenError("Failed to get stdin from ssh client. ", err)
		return
	}
	
	//在Session上起Shell
	mode := ssh.TerminalModes{
		ssh.ECHO : 0,
		ssh.TTY_OP_ISPEED : 14400,
		ssh.TTY_OP_OSPEED : 14400,
	}
	if err := sshe.session.RequestPty("xterm", 80, 40, mode); err != nil {
		sshe.clearSessionAndClientWhenError("Request for pseudo terminal failed. ", err)
	    return
	}
	err = sshe.session.Shell()
	if err != nil {
		sshe.clearSessionAndClientWhenError("Failed to start shell: ", err)
		return
	}
	
	// 修改提示符
	stdin.Write([]byte("export PS1='\\e\\a\\e\\a\\e'\n"))
	
	// 获取提示符, 如果10s没有收到, 则超时
	timeout := make(chan int, 1)
	go func(){
		time.Sleep(time.Second * 10)
		timeout <- 1
	}()
	select {
		case <-prompt_notify:
		// do nothing
		case <-timeout:
			sshe.clearSessionAndClientWhenError("Shell not support, or start timeout. ", errors.New("Shell not support, or start timeout."))
			return
	}
	
	//输入登录信息到log
	log.Println("Login Infomation:")
	log.Println(loginInfoBuffer.GetContent())
	
	//将isLogin设置为true
	sshe.isLogin = true
}

//在出错时清理Session和Client
func (sshe *SshExecutor) clearSessionAndClientWhenError(msg string, err error) {
	log.Println("ERROR: "+msg, err)
	if sshe.session != nil {
		warn := sshe.session.Close()
		if warn != nil {
			log.Println("WARN:  Session Close Error", warn)
		}
		sshe.session = nil
	}
	if sshe.client != nil {
		warn := sshe.client.Close()
		if warn != nil {
			log.Println("WARN:  Client Close Error", warn)
		}
		sshe.client = nil
	}
}

//销毁执行器
func (sshe *SshExecutor) Destory() {
	if false == sshe.isLogin {
		//如果没有登录直接返回
		log.Println("ERROR: SSH Client not login.")
		return
	}
	//销毁session
	if sshe.session != nil {
		err := sshe.session.Close()
		if err != nil {
			log.Println("WARN:  Session Close error. ", err)
		}
	}
	//断开连接
	if sshe.client != nil {
		err := sshe.client.Close()
		if err != nil {
			log.Println("WARN:  SSH Client Close error. ", err)
		}
	}
	sshe.isLogin = false
}

// 执行指令, 执行的指令为cmd, 如果timeout秒没有执行完成, 将会超时, 返回错误
func (sshe *SshExecutor) Execute(cmd string, timeout int) ([]ResultOutput, error){
	//当SSH客户端没有登录成功时, 直接返回错误
	if false == sshe.IsLogin(){
		return nil, errors.New("Client not login.")
	}
	//TODO 重定向输出
	//TODO 等待指令完成或超时
	//TODO 从输出中获取结果
	return nil, nil
}

//SshExcutor是否完成了登录
func (sshe *SshExecutor) IsLogin() bool {
	return sshe.isLogin
}

//singleRawWriterBuffer同一时间只允许一个携程进行写入, 同时不会对内容进行加工
type singleRawWriterBuffer struct {
	b  bytes.Buffer
	mu sync.Mutex
} 

//singleRawWriterBuffer的Write实现
func (w *singleRawWriterBuffer) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.b.Write(p)
}

//从singleRawWriterBuffer中获取内容, 内容一旦获取, 将从Buffer中删除
func (w *singleRawWriterBuffer) GetContent() string {
	w.mu.Lock()
	defer w.mu.Unlock()
	content := w.b.String()
	w.b.Reset()
	return content
}

// 一个拥有通知功能的Writer装饰折类, 当出现有系统设置的提示符时, 会通过ch这个
// 通道发送通知
type decoratorWriterForNofityer struct {
	w io.Writer
	cache []byte
	ch chan int
	stopOutput bool
}

// 装饰者的Write方法, 这个方法会将特定的提示符过滤不输出, 并在特定提示符出现时,
// 给通道发送通知
func (w *decoratorWriterForNofityer) Write(p []byte) (n int, err error){
	w.cache = []byte(string(w.cache) + string(p))
	// 从之前的输出中找特定提示符
	inx := bytes.Index(w.cache, promptings)
	if inx > -1 {
		//如果找到提示符, 通过ch进行通知
		w.ch <- 1
		w.cache = make ([]byte, 0)
	}
	// 输出内容, 但'\033'开始一直到特定提示符之前的文字, 不进行输出
	// ( '\033' 和 特定提示符有可能不在同一次调用中出现, 这样会有一些脏字符输出,
	// 但由于可能性小, 并且在正常输出之后, 这里忽略这种复杂情况)
	inx033 := bytes.IndexByte(p, 27)
	if(inx033 > -1){
		out := p[0: inx033]
		_, err = w.w.Write(out)
		inxPrompt := bytes.Index(p, promptings)
		if(inxPrompt > -1){
			out = p[inxPrompt + len(promptings):]
			_, err = w.w.Write(out)
		}
	} else {
		_, err = w.w.Write(p)
	}
	return len(p), err
}

// 生成一个写装饰者通知类
func newDecoratorWriterForNofityer(wi io.Writer, ch chan int) *decoratorWriterForNofityer{
	ret := decoratorWriterForNofityer{w:wi}
	ret.cache = make([]byte, 0)
	ret.ch = ch
	return &ret
}

type ResultOutput interface {
	Type() string
	String() string
}

type StderrOutput struct {
	content string
}

func (out StderrOutput) Type() string{
	return "ERROR"
}

func (out StderrOutput) String() string{
	return out.content
}

type StdoutOutput struct {
	content string
}

func (out StdoutOutput) Type() string{
	return "INFO"
}

func (out StdoutOutput) String() string{
	return out.content
}

