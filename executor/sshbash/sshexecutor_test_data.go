package sshbash

import (

)

var rightPasswordLoginInfo = SSHLoginInfo {
			Ip: "ip",   //change here
			Username: "username", //change here
			Password: "password", //change here
		}

var wrongPrivateKeyLoginInfo = SSHLoginInfo {
			Ip: "ip",   //change here
			Username: "username", //change here
			Key: []byte("abcd"),	//这里是错误的key
		}

