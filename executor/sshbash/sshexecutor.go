package sshbash

import (
	"bytes"
	"errors"
	"github.com/ratstars/ops-deployer/executor/commons"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"strings"
	"sync"
	"time"
)

var promptings = []byte{
	27, 7, 27, 7, 27,
}

// SSH的登录信息, Key和string只能有一个生效, Key如果不为nil, 将优先使用
type SSHLoginInfo struct {
	Ip       string
	Port     string
	Username string
	Password string
	Key      []byte
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
			Auth: []ssh.AuthMethod{
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
	LoginInfo     *SSHLoginInfo
	client        *ssh.Client
	session       *ssh.Session
	stdin         io.WriteCloser
	buff          MixWriterBuffer
	prompt_notify chan int
}

// 输入输出类型, 为OutputBuffer.buffByteType所用
const (
	UNKNOWN = 0
	STDOUT  = 1
	STDERR  = 2
)

//一次输入的片段
type OutputBufferPiece struct {
	buffType int
	content  []byte
}

//创建一个新的
func NewOutputBufferPiece(buffType int, content []byte) OutputBufferPiece {
	c := make([]byte, len(content))
	for i, v := range content {
		c[i] = v
	}
	return OutputBufferPiece{
		buffType: buffType,
		content:  c,
	}
}

//输入穿冲区,这个穿冲区的读写都会多线程安全的
type OutputBuffer struct {
	mu          sync.Mutex
	contentBuff []OutputBufferPiece
}

//写入一个缓冲区数据
func (buff *OutputBuffer) Write(piece OutputBufferPiece) {
	buff.mu.Lock()
	defer buff.mu.Unlock()
	buff.contentBuff = append(buff.contentBuff, piece)
}

//重组缓冲区片段, 如果连续的片段类型相同, 它将合成一个大的片段, 以保证不出现意外的换行
//这种换行在真实执行指令时是不存在的
func reorgBufferPiece(ori []OutputBufferPiece) []OutputBufferPiece {
	if len(ori) <= 1 {
		return ori
	}
	//要返回的空slice
	result := make([]OutputBufferPiece, 0)
	//将每一个取出
	tmp := ori[0]
	//查看后边每一个(当前值)是否与tmp相同类型
	for i := 1; i < len(ori); i++ {
		p := ori[i]
		if p.buffType == tmp.buffType {
			//如果内容相同, 将当前值加入tmp的内容中
			tmp.content = append(tmp.content, p.content...)
		} else {
			//否则将tmp加入到result后,当前值设为tmp
			result = append(result, tmp)
			tmp = p
		}
	}
	//将tmp加入result
	result = append(result, tmp)
	return result
}

//获取输出结果, 并将返回清空
func (buff *OutputBuffer) GetOutputAndClear() []commons.ResultOutput {
	result := make([]commons.ResultOutput, 0, 20)
	buff.mu.Lock()
	defer buff.mu.Unlock()
	// 对输出进行紧排
	buff.contentBuff = reorgBufferPiece(buff.contentBuff)
	// TODO 将contentBuff组合成ResultOutput返回

	return result
}

// ==========================TODO====================================

//初始化SSH执行器
func (sshe *SshExecutor) Init() {
	//如果已经初始化, 只是返回
	if sshe.isLogin {
		return
	}
	//检查各参数是否已经OK，如果有成员不完整用默认代替或提示出错
	if nil == sshe.LoginInfo {
		log.Println("ERROR: SSH Login Information is empty. Ssh is not start.")
		return
	}

	//登录
	config, err := sshe.LoginInfo.getClientConfig()
	if err != nil {
		log.Println("ERROR: get Client Config error, errorMsg: ", err)
		return
	}
	var address string
	if sshe.LoginInfo.Port == "" {
		address = sshe.LoginInfo.Ip + ":22"
	} else {
		address = sshe.LoginInfo.Ip + ":" + sshe.LoginInfo.Port
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

	//重定向输入输出
	sshe.prompt_notify = make(chan int)
	sshe.session.Stdout = newDecoratorWriterForNofityer(
		&BufferWriter{
			buff:       &sshe.buff,
			resultFunc: commons.NewStdoutOutput,
		}, sshe.prompt_notify)
	sshe.session.Stderr = &BufferWriter{
		buff:       &sshe.buff,
		resultFunc: commons.NewStdoutOutput,
	}
	sshe.stdin, err = sshe.session.StdinPipe()
	if err != nil {
		sshe.clearSessionAndClientWhenError("Failed to get stdin from ssh client. ", err)
		return
	}

	//在Session上起Shell
	mode := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
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

	// 启动bash, 并修改提示符
	sshe.stdin.Write([]byte("bash\n"))
	sshe.stdin.Write([]byte("export PS1='\\e\\a\\e\\a\\e'\n"))

	// 获取提示符, 如果10s没有收到, 则超时
	timeout := make(chan int, 1)
	go func() {
		time.Sleep(time.Second * 10)
		timeout <- 1
	}()
	select {
	case <-sshe.prompt_notify:
	// do nothing
	case <-timeout:
		sshe.clearSessionAndClientWhenError("Shell not support, or start timeout. ", errors.New("Shell not support, or start timeout."))
		return
	}

	//输入登录信息到log
	log.Println("Login Infomation:")
	log.Println(sshe.buff.GetOutputSetAndClear())

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
	//清空buff和关闭通道
	sshe.buff.GetOutputSetAndClear()
	close(sshe.prompt_notify)
	sshe.prompt_notify = nil
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
func (sshe *SshExecutor) Execute(cmd string, timeout int) ([]commons.ResultOutput, error) {
	//当SSH客户端没有登录成功时, 直接返回错误
	if false == sshe.IsLogin() {
		return nil, errors.New("Client not login.")
	}

	//等待指令完成或超时
	sshe.stdin.Write([]byte(cmd + "\n"))
	timeout_ch := make(chan int, 1)
	go func() {
		var st int64 = int64(time.Second)
		st = st * int64(timeout)
		time.Sleep(time.Duration(st))
		timeout_ch <- 1
	}()
	select {
	case <-sshe.prompt_notify:
		// do nothing
	case <-timeout_ch:
		err := errors.New("Execute cmd timeout. ")
		return nil, err
	}
	//从输出中获取结果
	return sshe.buff.GetOutputSetAndClear(), nil
}

//SshExcutor是否完成了登录
func (sshe *SshExecutor) IsLogin() bool {
	return sshe.isLogin
}

// 和IsLogin相同, 但这个方法是executor.Executor接口的一个方法
func (sshe *SshExecutor) IsReady() bool {
	return sshe.IsLogin()
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
	w          io.Writer
	cache      []byte
	ch         chan int
	stopOutput bool
}

// 装饰者的Write方法, 这个方法会将特定的提示符过滤不输出, 并在特定提示符出现时,
// 给通道发送通知
func (w *decoratorWriterForNofityer) Write(p []byte) (n int, err error) {
	//是否提示已经完成
	needNotify := false
	w.cache = []byte(string(w.cache) + string(p))
	// 从之前的输出中找特定提示符
	inx := bytes.Index(w.cache, promptings)
	if inx > -1 {
		//如果找到提示符, 通过ch进行通知
		needNotify = true
		w.cache = make([]byte, 0)
	}
	// 输出内容, 但'\033'开始一直到特定提示符之前的文字, 不进行输出
	// ( '\033' 和 特定提示符有可能不在同一次调用中出现, 这样会有一些脏字符输出,
	// 但由于可能性小, 并且在正常输出之后, 这里忽略这种复杂情况)
	inx033 := bytes.IndexByte(p, 27)
	if inx033 > -1 {
		out := p[0:inx033]
		_, err = w.w.Write(out)
		inxPrompt := bytes.Index(p, promptings)
		if inxPrompt > -1 {
			out = p[inxPrompt+len(promptings):]
			_, err = w.w.Write(out)
		}
	} else {
		_, err = w.w.Write(p)
	}
	if true == needNotify {
		w.ch <- 1
	}
	return len(p), err
}

// 生成一个写装饰者通知类
func newDecoratorWriterForNofityer(wi io.Writer, ch chan int) *decoratorWriterForNofityer {
	ret := decoratorWriterForNofityer{w: wi}
	ret.cache = make([]byte, 0)
	ret.ch = ch
	return &ret
}

// 混合写缓存, 这个类会将STDOUT和STDERR混合成一个BUFFER,并保证顺序. 同时结果保存
// 为[]ResultOutput, 方便后续调用者
type MixWriterBuffer struct {
	//锁
	mu sync.Mutex
	//上一次输入是否输出完一行
	lastLineFinish bool
	//存储
	outputSlice []commons.ResultOutput
}

//将ResultOutput加入到buffer中, 如果上一行没有完isLineFinished 为false, 这样如
//果下一次加入的ResultOutput类型与上一次的相同, 会将两次的类型合并
func (buff *MixWriterBuffer) Add(resultOutput commons.ResultOutput, isLineFinished bool) {
	buff.mu.Lock()
	defer buff.mu.Unlock()
	if nil == buff.outputSlice {
		buff.outputSlice = make([]commons.ResultOutput, 0, 20)
		buff.lastLineFinish = true
	}
	// 如果上一行没有完，考虑是否要是进行合并,
	if false == buff.lastLineFinish {
		l := len(buff.outputSlice)
		if l > 0 && buff.outputSlice[l-1].Type() == resultOutput.Type() {
			buff.outputSlice[l-1] = commons.Merge(buff.outputSlice[l-1], resultOutput)
			buff.lastLineFinish = isLineFinished
			return
		}
	}
	buff.outputSlice = append(buff.outputSlice, resultOutput)
	buff.lastLineFinish = isLineFinished
}

//获取Buffer的输出结果, 并清空
func (buff *MixWriterBuffer) GetOutputSetAndClear() []commons.ResultOutput {
	buff.mu.Lock()
	defer buff.mu.Unlock()
	buff.lastLineFinish = true
	ret := buff.outputSlice
	buff.outputSlice = make([]commons.ResultOutput, 0, 20)
	return ret
}

//输出的缓存写对象
type BufferWriter struct {
	buff       *MixWriterBuffer
	resultFunc func(s string) commons.ResultOutput
}

//BufferWriter的写方法
func (w *BufferWriter) Write(p []byte) (int, error) {
	sep := strings.Split(string(p), "\n")
	l := len(sep)
	var lastLineFinished bool
	//查看最后一行是否为空行, 如果是空行表示输出完成
	if sep[l-1] == "" {
		lastLineFinished = true
		//去除空行
		sep = sep[:l-1]
	}
	for i, line := range sep {
		if i < len(sep)-1 {
			// 非最后一行
			w.buff.Add(w.resultFunc(strings.TrimRight(line, "\r")), true)
		} else {
			w.buff.Add(w.resultFunc(strings.TrimRight(line, "\r")), lastLineFinished)
		}
	}
	return len(p), nil
}
