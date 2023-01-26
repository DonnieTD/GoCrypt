package main

import (
	"errors"
	"fmt"
	"os"
)

func Usage(err string) {
	if err != "" {
		fmt.Println("Error: " + err)
	}
	fmt.Println("Usage:")
	fmt.Println("	./gocrypt e ./PathToFolder Username")
	fmt.Println("	./gocrypt d ./PathToFile Username")
	if err != "" {
		os.Exit(1)
	}
}

func UsageParamsCheck(args []string) (string, string, string, string) {
	mode := args[1]
	PathToFileOrFolder := args[2]
	username := args[3]
	password := args[3]
	fmt.Println(os.Args)
	if len(args) != 5 {
		Usage("Wrong number of arguments")
	}

	if mode != "e" && mode != "d" {
		Usage("Mode does not exist")
	}

	if _, err := os.Stat(PathToFileOrFolder); errors.Is(err, os.ErrNotExist) {
		RaiseError("File/Folder does not exist: "+PathToFileOrFolder, true)
	}

	if len(username) < 5 {
		RaiseError("Username must be greater than 5 characters", true)
	}

	return mode, PathToFileOrFolder, username, password
}
