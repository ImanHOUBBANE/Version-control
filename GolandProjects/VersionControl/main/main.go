package main

import (
	"VersionControl/add"
	"VersionControl/checkout"
	"VersionControl/commit"
	"VersionControl/log"
	"VersionControl/login"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	if len(os.Args) < 2 {
		showHelp()
		return
	}

	command := os.Args[1]
	switch command {
	case "config":
		login.HandleConfigCommand()
	case "add":
		add.HandleAddCommand()
	case "commit":
		commit.HandleCommitCommand()
	case "log":
		log.HandleLogCommand()
	case "checkout":
		checkout.HandleCheckoutCommand()
	case "--help":
		showHelp()

	default:

		fmt.Printf("'%s' is not a SVCS command.", os.Args[1])
	}
}

func showHelp() {
	fmt.Println("These are SVCS commands:")
	fmt.Println("config     Get and set a personal information.")
	fmt.Println("add        Add a file to the index.")
	fmt.Println("log        Show commit logs.")
	fmt.Println("commit     Save changes.")
	fmt.Println("checkout   Restore a file.")
}
