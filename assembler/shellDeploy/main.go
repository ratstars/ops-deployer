package main

import (
	"flag"
	"fmt"
	"github.com/ratstars/ops-deployer/controller"
	"github.com/ratstars/ops-deployer/script"
	"github.com/ratstars/ops-deployer/view"
	"io/ioutil"
	"os"
)

const APP_VERSION = "0.1"

// The flag package provides a default help printer via -h switch
var versionFlag *bool = flag.Bool("ver", false, "Print the version number.")

var scriptFlag *string = flag.String("s", "", "Input Script File Name.")

func main() {
	flag.Parse() // Scan the arguments list

	if *versionFlag {
		fmt.Println("Version:", APP_VERSION)
		return
	}

	if "" == *scriptFlag {
		flag.Usage()
		return
	}

	//1. open script file
	bytes, err := ioutil.ReadFile(*scriptFlag)
	if err != nil {
		fmt.Errorf("File Open Failed. %v\n", err)
		os.Exit(-1)
		return
	}
	content := string(bytes)

	//2. create Decoder and decode script texts
	decoder := &script.Decoder{}
	script, err := decoder.Decode(content)
	if err != nil {
		fmt.Errorf("Decode Script Files Error. %v\n", err)
		os.Exit(-1)
		return
	}

	//3. create View
	view := &view.ShellView{}
	//4. create Controller
	ctrler := controller.DefaultController{
		Confirmer: view,
		View:      view,
	}
	//5. run script and print result.
	err = ctrler.RunScript(script)
	if err != nil {
		fmt.Errorf("Run Script Error. %v\n", err)
	}
	fmt.Println("Script Execution Finished. ")
}
