package sshbash

import (
    "testing"
)

// 测试IsLogin是否可以正常返回
func TestIsLogin(t *testing.T) {
	ssh := &SshExecutor {
		LoginInfo: &SSHLoginInfo{
			Ip:"localhost",
			Username: "username",
			Password: "password",
		},
	}
	if ssh.IsLogin() != false {
		t.Error("Excpect ssh.IsLogin() is false")
	}
}

// 测试如果输入一个空的LoginInfo的行为
func TestInitEmptyLoginInfo(t *testing.T) {
	ssh := &SshExecutor {}
	ssh.Init()
	if ssh.IsLogin() != false {
		t.Error("Not give Login info to ssh, expect not login status after Init()")
	}
}

// 测试用密码登录的情况
func TestLoginWithPassword(t *testing.T) {
	ssh := &SshExecutor {
		LoginInfo: &rightPasswordLoginInfo,
	}
	ssh.Init()
	if false == ssh.IsLogin() {
		t.Error("Expect Login Successful, but get failed result.")
	}
	ssh.Destory()
}

// 测试用秘钥登录的情况
func TestLoginWithKey(t *testing.T) {
	ssh := &SshExecutor {
		LoginInfo: &wrongPrivateKeyLoginInfo,
	}
	ssh.Init()
	if true == ssh.IsLogin() {
		t.Error("Wrong Key expect not login result.")
		return
	}
	keyS := "-----BEGIN RSA PRIVATE KEY-----\n"+
		"MIIEowIBAAKCAQEA3ReQn3mBz3X6bn216vbaTxZfU4KWuqILcD3suHLpTAaw5OFf\n"+
		"O6Dg0kp7WVAAoxiFUl7/rXpRGK6MbMD/NZB0vwuButs1y+4ZK8oTR4lcXgZUOfIJ\n"+
		"lgInzF/SF5b0m2rPMXMHRa/JQ3SUmug7NquA4TZbgbAj9klF6qplxiPyBt6pXvMr\n"+
		"crZu+28OAHo3oU94oOWzoPB46Hw4hJxzk2Y/fnso4Gznn+AegThAS4/xyzb6yavc\n"+
		"1AvhoRAn/7muXP+ZH3mr8w99QTMZTJfpQNr3sApqEQuq2HyD/7m81xq+D95Kay3W\n"+
		"FEpJGDxJztrs1NLnITQOKqMluPRItnBQJNlM4wIBIwKCAQEAl5shV2lDEep/1rya\n"+
		"ADQ97RanxDxKGZOwEnObAiLpHjB5TH1Ineqoyra6+2oPEMBba67bNSCsozXcolh0\n"+
		"fIBQDflDA8mD+Y1TFrZznsSXG+cVLwxeWDv+CHw4SrCnuwdpgP5rYvwyPOI6A9Jx\n"+
		"vxaEQqjuSkzlda6WV8VNGidG4CJhCmcxbPScgou8xojFA8b842aImPRFfHTHWMYZ\n"+
		"jkMaVwtxFoBK8zZ6rtgT0s1RXkNNx7FF2F9b6jIvom8dXD4I2MLt2lbwcFFnzyuW\n"+
		"SCPXg0nqWWUbiB756fgM7oZYqUaC9o9xv0QLaA2++GoyoKuFSkapZaqnQDK7gma+\n"+
		"l9HcSwKBgQDxA1ywUP/b/S1Zdjnnpe8AWI8sh6CZxhFDw941LUq72/LKtlkJERC/\n"+
		"LjFFufbUp/7Mdl9/xJf64q7yYbQ9m9nBPwegnSnsfUuxCY6FI+ah3Y8l2Pakv44P\n"+
		"moEoQrAEst3hhshcha6lgiHEwt+F9g1+4ltzupJaOmQCa4j0PI/exQKBgQDq1xVT\n"+
		"DXyEFTXr6LjbVd5HfLZVOkD5zOtSCfKudqRSjkJZPpjGlsqXFAvdqNhbzxYHkBSB\n"+
		"3pOAs4KwJni7RUtgTwBwd3wcysYGNOnTnmq7iFojyww98mvymq7XDebNZiQ/PNRJ\n"+
		"AlhoSQZ5NwsyXyCbMobjZnfs65ufghGF9Yy3hwKBgQCeYUuJztQFl763Ic5HxNBC\n"+
		"Dk91CKtdvKxCeWYi8eCnVgXy7NtsW6vrWN6M5+tYi6dwawuOeeA3Jz/D2c43HUX0\n"+
		"BNkgZ0dvhYmC96bMhU5qXmViA5rECNmye3lyOnOrUPg1Hg6jM0bhyos4KEm+bn3l\n"+
		"qrEgKiWpAcyxIhgreEFJPwKBgA1rYE3jg3VDCmVAf5eBP+bT7SlxCwb1xE3Ursgk\n"+
		"CWPNnWQvdnG/eUp2LJBSyojnQxZgASv+Fw6rLAoQ07LuBE6lbb1IqAGlL+MY92st\n"+
		"n7Lx2UPft46C4ZjVo5dCn3lzjQrtiHkzVYJNUNO6AKPK6+ucfLyJgzIcF4V1JZKg\n"+
		"US8PAoGBAIFEdjML8UD0zrvWGRgki+6ALBoWZxInn4vmvfVq36oc6bTcKCfBcCPW\n"+
		"kRvvyMbOCkMST6UI7kl58Vso7/UiQG+qy1+ZJet5P3kt0rPGMVB1K5/FrKNO/SHw\n"+
		"r6Q6KVxUIZtc/AFr2Rq7GIrxwSNdy+oIgu78z/FsoFO4htJP5See\n"+
		"-----END RSA PRIVATE KEY-----"
	ssh.LoginInfo.Key = []byte(keyS)	//这里是正确的key
	ssh.Init()
	if false == ssh.IsLogin() {
		t.Error("Expect Login Successful, but get failed result.")
	}
	ssh.Destory()
}

//测试销毁
func TestDestory(t *testing.T) {
	ssh := &SshExecutor {
		LoginInfo: &rightPasswordLoginInfo,
	}
	ssh.Init()
	if false == ssh.IsLogin() {
		t.Error("Expect Login Successful, but get failed result.")
		return
	}
	ssh.Destory()
	if true == ssh.IsLogin() {
		t.Error("Expect Logout Successful, but get IsLogin() == true. ")
		return
	}
	ssh.Init()
	if false == ssh.IsLogin() {
		t.Error("Expect Login Successful Again, but get failed result.")
		return
	}
	ssh.Destory()
	if true == ssh.IsLogin() {
		t.Error("Expect Logout Successful Again, but get IsLogin() == true. ")
		return
	}
}

func TestExecute(t *testing.T) {
	ssh := &SshExecutor {
		LoginInfo: &rightPasswordLoginInfo,
	}
	ssh.Init()
	if false == ssh.IsLogin() {
		t.Error("Expect Login Successful, but get failed result.")
		return
	}
	resultSet, err := ssh.Execute("echo ABCD", 1)
	if err != nil {
		t.Error("Error Occures: ", err)
		return
	}
	if len(resultSet) != 1 {
		t.Error("Expect One Line Return")
		return
	}
}
