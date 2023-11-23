package main

import (
	"fmt"
	"runtime"
)

var (
	description    = "Cloud Carbon collection exporter"
	gitSHA         = "n/a"
	name           = "Cloud Carbon"
	source         = "https://github.com/re-cinq/cloud-carbon"
	version        = "0.0.1-dev"
	refType        = "branch" // branch or tag
	refName        = ""       // the name of the branch or tag
	buildTimestamp = ""
)

func PrintVersion() {
	fmt.Printf("Name:           %s\n", name)
	fmt.Printf("Version:        %s\n", version)
	fmt.Printf("RefType:        %s\n", refType)
	fmt.Printf("RefName:        %s\n", refName)
	fmt.Printf("Git Commit:     %s\n", gitSHA)
	fmt.Printf("Description:    %s\n", description)
	fmt.Printf("Go Version:     %s\n", runtime.Version())
	fmt.Printf("OS / Arch:      %s / %s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("Source:         %s\n", source)
	fmt.Printf("Built:          %s\n", buildTimestamp)

}
