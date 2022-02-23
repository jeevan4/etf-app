package main

import (
	"etf-app/cli"
	"fmt"
	"os"
	"time"
)

func main() {
	fmt.Println("Welcome to the ETF App")
	start := time.Now()
	args := os.Args
	etfCommand := cli.NewEtfCommand("Etf CLI App")
	if len(args) < 2 {
		// fmt.Printf("Not Enough arguments privided\nPlease provide\n\tetf <command> --help\n")
		etfCommand.Help()
		return
	}
	_, ok := etfCommand.Commands[args[1]]
	if ok {
		etfCommand.Commands[args[1]](args[2:])
	} else {
		fmt.Printf("%s not recognized! Please follow help\n", args[1])
		etfCommand.Help()
	}

	// fmt.Printf("Args %s\n", args[1:])
	fmt.Println("Command Elapsed:", time.Since(start))
}
