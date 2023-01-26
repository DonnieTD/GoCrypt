package main

import (
	"fmt"
	"os"
)

func RaiseError(err string, shouldExit bool) {
	fmt.Println("Error: " + err)
	if shouldExit {
		os.Exit(1)
	}
}
