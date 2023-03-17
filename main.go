package main

import (
	"fmt"
	"gdos/modules/core"
	"gdos/modules/methods"
)

func main() {
	core.ClearConsole()
	fmt.Println("Enter an target IPV4 address.")
	fmt.Println("")
	var target string
	fmt.Scanln(&target)
	core.ClearConsole()
	fmt.Println("Enter how many seconds it should last.")
	fmt.Println("")
	var last string
	fmt.Scanln(&last)

	methods.IMCP(target, last)

}
