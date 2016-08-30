package main

import (
	"flag"
	"fmt"
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
}
