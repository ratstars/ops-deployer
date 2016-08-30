package script

import (
    "testing"
)

func TestNoEmptySplitN(t *testing.T) {
	in := "define executor1 ssh {\"ip\":\"ip\", \"username\":\"username\", \"password\":\"password\"}"
	out := NoEmptySplitN(in, 4)
	if len(out) != 4 {
		t.Error("Expect 4 part")
	}
	in = "  define executor2   ssh {\"ip\":\"ip\", \"username\":\"username\", \"password\":\"password\"}  "
	out = NoEmptySplitN(in, 4)
	if len(out) != 4 {
		t.Error("Expect 4 part")
	}
}

func TestDecoderExecutor1(t *testing.T) {
	d := Decoder{}
	script1, err:= d.Decode("define executor1 ssh {\"ip\":\"ip\", \"username\":\"username\", \"password\":\"password\"}")
	if err != nil {
		t.Error("Decode Script1 Error: ", err)
		return
	}
	if len(script1.Executors) != 1 || len(script1.Commands) != 0 {
		t.Error("Exepct script1 only contains one executor, but len(executor), len(command) = ", len(script1.Executors),len(script1.Commands))
		return
	}
	if script1.Executors[0].Name != "executor1" || script1.Executors[0].Type != "ssh" {
		t.Error("Unexpect decode result. ", script1.Executors[0])
		return
	}
}

func TestDecoderExecutor1Error(t *testing.T) {
	d := Decoder{}
	_, err:= d.Decode("define executor1 ssh {\"ip\":\"ip, \"username\":\"username\", \"password\":\"password\"}")
	if err == nil {
		t.Error("Expect error.")
		return
	}

}

func TestDecoderExecutor2(t *testing.T) {
	input := "define executor1 ssh {\"ip\":\"ip\", \"username\":\"username\", \"password\":\"password\"}\n" + 
			"  define executor2   ssh {\"ip\":\"ip\", \"username\":\"username\", \"password\":\"password\"}  \n"
	
	d := Decoder{}
	script1, err:= d.Decode(input)
	if err != nil {
		t.Error("Decode Script1 Error: ", err)
		return
	}
	if len(script1.Executors) != 2 || len(script1.Commands) != 0 {
		t.Error("Exepct script1 only contains one executor, but len(executor), len(command) = ", len(script1.Executors),len(script1.Commands))
		return
	}
	if script1.Executors[1].Name != "executor2" || script1.Executors[0].Type != "ssh" {
		t.Error("Unexpect decode result. ", script1.Executors[0])
		return
	}
}

func TestDecoderExecutor1and1(t *testing.T) {
	input := "define executor1 ssh {\"ip\":\"ip\", \"username\":\"username\", \"password\":\"password\"}\n" + 
			"define executor1 ssh {\"ip\":\"ip\", \"username\":\"username\", \"password\":\"password\"}\n"
	
	d := Decoder{}
	_, err:= d.Decode(input)
	if err == nil {
		t.Error("Expect Dumplication Error.")
		return
	}

}

func TestDecodeComment1(t *testing.T) {
	input := "define executor1 ssh {\"ip\":\"ip\", \"username\":\"username\", \"password\":\"password\"}\n" +
			"   # this is a comment"
	d := Decoder{}
	script1, err := d.Decode(input)
	if err != nil {
		t.Error("Decode Error: ",err)
	}
	if len(script1.Commands) != 1 || len(script1.Executors) != 1 {
		t.Error("Expect 1 Commands and 1 Executor.")
		return
	}
	if script1.Commands[0].IsComment == false || script1.Commands[0].Command != "this is a comment" {
		t.Error("Error Comment: ", script1.Commands[0])
		return
	}
}

func TestDecodeCommand1(t *testing.T) {
		input := "define executor1 ssh {\"ip\":\"ip\", \"username\":\"username\", \"password\":\"password\"}\n" +
			"   executor1 echo ABCD\n" +
			"in 60  "
	d := Decoder{}
	script1, err := d.Decode(input)
	if err != nil {
		t.Error("Decode Error: ", err)
		return
	}
	if len(script1.Commands) != 1 || len(script1.Executors) != 1 {
		t.Error("Expect 1 Commands and 1 Executor.")
		return
	}
	if script1.Commands[0].ExecutorName != "executor1" || script1.Commands[0].Command != "echo ABCD" {
		t.Error("Executor1 not expect one")
		return
	}
}

func TestDecodeCommand2(t *testing.T) {
	input := "define executor1 ssh {\"ip\":\"ip\", \"username\":\"username\", \"password\":\"password\"}\n" +
			"   executor1 echo ABCD\n" +
			"in 60  expect"
	d := Decoder{}
	script1, err := d.Decode(input)
	if err == nil {
		t.Error("Expect Error", script1)
		return
	}
}

func TestDecodeCommand3(t *testing.T) {
	input := "define executor1 ssh {\"ip\":\"ip\", \"username\":\"username\", \"password\":\"password\"}\n" +
			"   executor1 echo ABCD\n" +
			"in 60  expect ABCD"
	d := Decoder{}
	script1, err := d.Decode(input)
	if err != nil {
		t.Error("Decode Error: ", err)
		return
	}
	if len(script1.Commands) != 1 || len(script1.Executors) != 1 {
		t.Error("Expect 1 Commands and 1 Executor.")
		return
	}
	if script1.Commands[0].ExecutorName != "executor1" || script1.Commands[0].Command != "echo ABCD" ||
		script1.Commands[0].ExpectRegular != "ABCD"{
		t.Error("Executor1 not expect one")
		return
	}
}

func TestDecodeCommand4(t *testing.T) {
	input := "define executor1 ssh {\"ip\":\"ip\", \"username\":\"username\", \"password\":\"password\"}\n" +
			"   executor1 echo ABCD\n" +
			"in 60  \nexpect ABCD\n"
	d := Decoder{}
	_, err := d.Decode(input)
	script1, err := d.Decode(input)
	if err != nil {
		t.Error("Decode Error: ", err)
		return
	}
	if len(script1.Commands) != 1 || len(script1.Executors) != 1 {
		t.Error("Expect 1 Commands and 1 Executor.")
		return
	}
	if script1.Commands[0].ExecutorName != "executor1" || script1.Commands[0].Command != "echo ABCD" ||
		script1.Commands[0].ExpectRegular != "ABCD"{
		t.Error("Executor1 not expect one")
		return
	}
}

func TestDecodeCommand5(t *testing.T) {
	input := "define executor1 ssh {\"ip\":\"ip\", \"username\":\"username\", \"password\":\"password\"}\n" +
			"   executor1 echo ABCD\n" +
			"in 60  \nexpect ABCD\n" +
			"unexpect ACD"
	d := Decoder{}
	_, err := d.Decode(input)
	script1, err := d.Decode(input)
	if err != nil {
		t.Error("Decode Error: ", err)
		return
	}
	if len(script1.Commands) != 1 || len(script1.Executors) != 1 {
		t.Error("Expect 1 Commands and 1 Executor.")
		return
	}
	if script1.Commands[0].ExecutorName != "executor1" || script1.Commands[0].Command != "echo ABCD" ||
		script1.Commands[0].ExpectRegular != "ABCD"{
		t.Error("Executor1 not expect one")
		return
	}
}


