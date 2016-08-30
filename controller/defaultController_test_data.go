package controller

import ()

var script1_test = "define shell SSHBASH {\"ip\":\"ip\", \"username\":\"username\", \"password\":\"password\"}\n" +
	"shell echo ABCD\n" +
	"in 60  expect ABCD"

var script2_test = "define shell SSHBASH {\"ip\":\"ip\", \"username\":\"username\", \"password\":\"password\"}\n" +
	"shell echo ABCD\n" +
	"in 60  expect abcd"

var script3_test = "define shell1 SSHBASH {\"ip\":\"ip\", \"username\":\"username\", \"password\":\"password\"}\n" +
	"define shell2 SSHBASH {\"ip\":\"ip\", \"username\":\"username\", \"password\":\"password\"}\n" +
	"shell1 echo ABCD\n" +
	"in 60  expect ABCD\n" +
	"shell2 echo ABCD"

var script4_test = "define shell1 SSHBASH {\"ip\":\"ip\", \"username\":\"username\", \"password\":\"password\"}\n" +
	"define shell2 SSHBASH {\"ip\":\"ip\", \"username\":\"username\", \"password\":\"password\"}\n" +
	"shell1 echo ABCD\n" +
	"in 60  expect abcd\n" +
	"shell2 sleep 2\n" +
	"in 0"

var script5_test = "define shell SSHBASH {\"ip\":\"ip\", \"username\":\"username\", \"password\":\"password\"}\n" +
	"shell echo ABCD\n" +
	"in 60  unexpect ABCD"
