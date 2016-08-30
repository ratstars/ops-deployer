package controller

import (
	"github.com/ratstars/ops-deployer/script"
	"testing"
)

//测试用Confirmer, 确认结果提前输入好, 用于测试控制
type preEnterConfirmer bool

func (p *preEnterConfirmer) Confirm(info string) bool {
	return bool(*p)
}

func (p *preEnterConfirmer) DisplayAndPause(info string) {

}

func TestDefaultControllerRunScript1(t *testing.T) {
	d := script.Decoder{}
	script1, err := d.Decode(script1_test)
	if err != nil {
		t.Error("Decode Error: ", err)
		return
	}

	result := preEnterConfirmer(true)
	ctrler := &DefaultController{
		Confirmer: &result,
	}

	err = ctrler.RunScript(script1)
	if err != nil {
		t.Fatal("Run script error.", err)
	}
}

func TestDefaultControllerRunScript2(t *testing.T) {
	d := script.Decoder{}
	script2, err := d.Decode(script2_test)
	if err != nil {
		t.Error("Decode Error: ", err)
		return
	}

	result := preEnterConfirmer(true)
	ctrler := &DefaultController{
		Confirmer: &result,
	}

	err = ctrler.RunScript(script2)
	if err == nil {
		t.Fatal("Run script error. expect a error result")
	}
}

func TestDefaultControllerRunScript3(t *testing.T) {
	d := script.Decoder{}
	script3, err := d.Decode(script3_test)
	if err != nil {
		t.Error("Decode Error: ", err)
		return
	}

	result := preEnterConfirmer(true)
	ctrler := &DefaultController{
		Confirmer: &result,
	}

	err = ctrler.RunScript(script3)
	if err != nil {
		t.Fatal("Run script error.", err)
	}
}

func TestDefaultControllerRunScript4(t *testing.T) {
	d := script.Decoder{}
	script4, err := d.Decode(script4_test)
	if err != nil {
		t.Error("Decode Error: ", err)
		return
	}

	result := preEnterConfirmer(true)
	ctrler := &DefaultController{
		Confirmer: &result,
	}

	err = ctrler.RunScript(script4)
	if err == nil {
		t.Fatal("Run script error. expect a error result")
	}
}

func TestDefaultControllerRunScript5(t *testing.T) {
	d := script.Decoder{}
	script5, err := d.Decode(script5_test)
	if err != nil {
		t.Error("Decode Error: ", err)
		return
	}

	result := preEnterConfirmer(true)
	ctrler := &DefaultController{
		Confirmer: &result,
	}

	err = ctrler.RunScript(script5)
	if err == nil {
		t.Fatal("Run script error. expect a error result")
	}
}
