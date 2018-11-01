package main

import "github.com/gobuffalo/packr"

func printHelpMessage() {
	res := packr.NewBox("./res")
	help, error := res.FindString("help.txt")
	if (error != nil) {
		return
	}
	println(help)
}