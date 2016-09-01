package main

import (
	"flag"
	"fmt"
	"github.com/ratstars/ops-deployer/assembler"
	"github.com/ratstars/ops-deployer/view"
	"log"
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
	//1. change Log
	file, err := os.OpenFile("shellDeployer.log", os.O_CREATE|os.O_APPEND, 0)
	if err != nil {
		fmt.Errorf("Can not open log file. Log will output to stdout and stderr.")
	}
	defer func() {
		if file != nil {
			file.Close()
		}
	}()
	log.SetOutput(file)

	//2. open script
	script_file, err := os.Open(*scriptFlag)
	if err != nil {
		fmt.Errorf("Open Script File Error. %v", err)
		os.Exit(-1)
		return
	}
	defer script_file.Close()

	//3. create View and Confirmer
	view := &view.ShellView{}

	//4. create assembler and run
	deployer := &assembler.ShellDeployer{
		Confirmer: view,
		View:      view,
	}

	i := deployer.Run(script_file)
	os.Exit(i)
	return
}
