package executor

import (
	"encoding/json"
	"errors"
	"github.com/ratstars/ops-deployer/executor/sshbash"
)

type sshBashUserInfo struct {
	Ip       string `json:"ip"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func CreateExecutor(types string, args string) (Executor, error) {
	switch types {
	case "SSHBASH":
		return SSHBashExecutorCreator(args)
	default:
		return nil, errors.New("Not Support Executor's Type.")
	}
}

func SSHBashExecutorCreator(args string) (*sshbash.SshExecutor, error) {
	var sshInfo sshBashUserInfo
	err := json.Unmarshal([]byte(args), &sshInfo)
	if err != nil {
		return nil, err
	}
	ip := sshInfo.Ip
	username := sshInfo.Username
	password := sshInfo.Password
	if "" == ip {
		return nil, errors.New("SSHBASH Executor Args Error, Expect ip.")
	}
	if "" == username {
		return nil, errors.New("SSHBASH Executor Args Error, Expect username.")
	}
	if "" == password {
		return nil, errors.New("SSHBASH Executor Args Error, Expect password.")
	}
	ssh := &sshbash.SshExecutor{
		LoginInfo: &sshbash.SSHLoginInfo{
			Ip:       ip,
			Username: username,
			Password: password,
		},
	}
	return ssh, nil
}
