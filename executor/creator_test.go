package executor

import (
	"testing"
)

func TestSSHBashExecutorCreator(t *testing.T) {
	input := "{\"ip\":\"ip\", \"username\":\"username\", \"password\":\"password\"}"
	sshbash, err := SSHBashExecutorCreator(input)
	if err != nil {
		t.Fatal("Run SSHBashExecutorCreator error.", err)
	}
	if sshbash.LoginInfo.Ip != "ip" || sshbash.LoginInfo.Password != "password" || sshbash.LoginInfo.Username != "username" {
		t.Fatal("Not Expect Decoder Result. Result is:", sshbash)
	}
}
