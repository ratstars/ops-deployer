package script

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"log"
	"strings"
	"strconv"
)

// 脚本解析器
type Decoder struct {
}

// 解析输入的脚本
func (d *Decoder) Decode(input string) (Scripter, error) {
	return d.DecodeReader(strings.NewReader(input))
}

func (d *Decoder) DecodeReader(in io.Reader) (Scripter, error) {
	var scripter Scripter
	scripter.Executors = make([]ExecutorDescriber, 0, 5)
	scripter.Commands = make([]CommandDescriber, 0, 50)
	// 将输入解析成行
	scanner := bufio.NewScanner(in)
	lines := make([]string, 0, 50)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Println("DecodeReader read error: ", err)
		return Scripter{}, err
	}
	// 解析脚本, 解析后从lines中删除
	for len(lines) > 0 {
		// 解析define
		if strings.Index(strings.TrimSpace(lines[0]), "define ") == 0 {
			err := decodeExecutor(&lines, &scripter.Executors)
			if err != nil {
				log.Println("Decode Executor Error: ", err)
				return Scripter{}, err
			}
		} else if strings.TrimSpace(lines[0]) == "" {
			//空行, 删除空行
			lines = lines[1:]
		} else if strings.TrimSpace(lines[0])[0] == '#' {
			// 注释
			decodeComment(&lines, &scripter.Commands)
		} else {
			// 其它内容, 按指令进行解析
			err := decodeCommand(&lines, &scripter.Commands)
			if err != nil {
				log.Println("Decode Command Error: ", err)
				return Scripter{}, err
			}
		}
	}
	// 检查Executor的名字是否重复
	err := checkDumplicatedExecutorName(scripter.Executors)
	if err != nil {
		log.Println("Executors' Name Dumplicated.", err)
		return Scripter{}, err
	}
	// 检查cmd的执行器名字的有效性
	err = checkExecutorNameInCommand(&scripter) 
	if err != nil {
		log.Println("Command's executor error.", err)
		return Scripter{}, err
	}
	return scripter, nil
}

// 检查Executor是否有重复的名字
func checkDumplicatedExecutorName(executors []ExecutorDescriber) error {
	executorNameSet := make(map[string]bool)
	for _, v := range executors {
		if _, ok := executorNameSet[v.Name]; false == ok {
			executorNameSet[v.Name] = true
		} else {
			return errors.New("Executors' Name Dumplicated: "+v.Name)
		}
	}
	return nil
}

// 检查Command中的Executor的名字是否有效
func checkExecutorNameInCommand(script *Scripter) error {
	executorNameSet := make(map[string]bool)
	for _, v := range script.Executors {
		executorNameSet[v.Name] = true
	}
	for _, v := range script.Commands {
		if true == v.IsComment {
			continue
		}
		if _, ok := executorNameSet[v.ExecutorName]; false == ok {
			// 没有定义executor
			return errors.New("Executor Not Define: "+v.ExecutorName)
		}
	}
	return nil
}

type emptyStruct struct{}

// 检看name是否是合规的名字
func validName(name string) bool {
	if name == "expect" || name == "in" || name == "unexpect" || name == "define" {
		return false
	}
	return true
}

// 解析executor
func decodeExecutor(lines *[]string, executors *[]ExecutorDescriber) error {
	//已经确定是一个executor后调用这个函数
	line := (*lines)[0]
	parts := NoEmptySplitN(line, 4)
	if len(parts) != 4 {
		return errors.New("define line format error")
	}
	var executor ExecutorDescriber
	executor.Name = parts[1]
	executor.Type = parts[2]
	executor.Args = parts[3]
	// 检查Name的有效性
	if false == validName(executor.Name) {
		return errors.New("Invalid Executor Name")
	}
	// 检查Args是否是json格式
	i := &emptyStruct{}
	err := json.Unmarshal([]byte(executor.Args), i)
	if err != nil {
		return errors.New("Executor Args format error: " + err.Error())
	}
	(*executors) = append(*executors, executor)
	(*lines) = (*lines)[1:]
	return nil
}

func decodeCommand(lines *[]string, commands *[]CommandDescriber) error {
	//当不以define开头时，也不是空行时调用这个函数
	cmd := CommandDescriber{}
	line := strings.TrimSpace((*lines)[0])
	parts := NoEmptySplitN(line, 2)
	if len(parts) != 2 {
		return errors.New("Command line format error")
	}
	cmd.ExecutorName = parts[0]
	cmd.Command = parts[1]
	(*lines) = (*lines)[1:]
	//查看是否有in expect unexpect表达式
	for len(*lines) > 0 {
		line = (*lines)[0]
		if strings.Index(line, "in ") != 0 &&
			strings.Index(line, "expect ") != 0 &&
			strings.Index(line, "unexpect ") != 0 {
			//不是in expect unexpect就退出
			break
		}
		if strings.Index(line, "in") == 0 {
			parts = NoEmptySplitN(line, 3)
			if len(parts) < 2 {
				return errors.New("in statement format error, expect a number.")
			}
			timeout, err := strconv.Atoi(parts[1])
			if err != nil {
				log.Print("in statement format error:", err)
				return err
			}
			cmd.Timeout = timeout
			if len(parts) >= 3 {
				// 如果in表达式后还有内容
				// 判断表达式是否合法
				if strings.Index(parts[2], "expect ") != 0 && strings.Index(parts[2], "unexpect ") != 0 {
					//只支持expect或unexpect
					return errors.New("in statement format error, unknown statements.")
				}
				(*lines)[0] = parts[2]
			} else {
				(*lines) = (*lines)[1:]
			}
		}
		if strings.Index(line, "expect ") == 0 {
			parts = NoEmptySplitN(line, 2)
			if len(parts) != 2 {
				return errors.New("expect statement format error, expect a regular expression.")
			}
			cmd.ExpectRegular = parts[1]
			(*lines) = (*lines)[1:]
		}
		if strings.Index(line, "unexpect") == 0 {
			parts = NoEmptySplitN(line, 2)
			if len(parts) != 2 {
				return errors.New("unexpect statement format error, expect a regular expression.")
			}
			cmd.UnexpectRegular = parts[1]
			(*lines) = (*lines)[1:]
		}
	}
	*commands = append(*commands, cmd)
	return nil
}

func decodeComment(lines *[]string, commands *[]CommandDescriber) {
	line := (*lines)[0]
	comment := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), "#"))
	cmd := CommandDescriber {
		IsComment: true,
		Command: comment,
	}
	*lines = (*lines)[1:]
	*commands = append(*commands, cmd)
}

// 按sep切分s，但如果切分出的子串有空串，会删除空串
// 对于res小于等于0的情况，会返回nil而不是空串
func NoEmptySplitN(s string, res int) [] string {
	if res <= 0 {
		return nil
	}
	ss := strings.TrimSpace(s)
	//空串返回空的slice
	slice := make([]string, 0, res)
	if ss == "" {
		return slice
	}
	if res <= 1 {
		slice = append(slice, ss)
		return slice
	}
	part := strings.SplitN(ss, " ", 2)
	if len(part) == 1 {
		slice = append(slice, ss)
		return slice
	}
	slice = append(slice, part[0])
	slice = append(slice, NoEmptySplitN(strings.TrimSpace(part[1]), res -1)...)
	return slice
}
