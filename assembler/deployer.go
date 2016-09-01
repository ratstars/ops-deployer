package assembler

import (
	"github.com/ratstars/ops-deployer/controller"
	"github.com/ratstars/ops-deployer/script"
	"github.com/ratstars/ops-deployer/view"
	"io"
	"io/ioutil"
)

type ShellDeployer struct {
	//确认器
	Confirmer view.Confirmer
	//视图
	View view.View
}

// 部署器的执行, 系统的输入参数是脚本内容的Reader
// 返回值为整数, 0表示成功, 非0表示失败, 这个值可以作为程序的返回值
func (d *ShellDeployer) Run(scriptReader io.Reader) int {
	//1. read script reader
	bytes, err := ioutil.ReadAll(scriptReader)
	if err != nil {
		d.View.DisplayInfo("Script Read Failed.")
		d.View.DisplayInfo(err.Error())
		return -1
	}
	content := string(bytes)

	//2. create Decoder and decode script texts
	decoder := &script.Decoder{}
	script, err := decoder.Decode(content)
	if err != nil {
		d.View.DisplayInfo("Decode Script Error.")
		d.View.DisplayInfo(err.Error())
		return -1
	}

	//3. create Controller
	ctrler := controller.DefaultController{
		Confirmer: d.Confirmer,
		View:      d.View,
	}
	//4. run script and print result.
	err = ctrler.RunScript(script)
	if err != nil {
		d.View.DisplayInfo("Run Script Error")
		d.View.DisplayInfo(err.Error())
	}
	d.View.DisplayInfo("Script Execution Finished.")
	return 0
}
