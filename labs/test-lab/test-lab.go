package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Please enter an argument")
		os.Exit(0)
	} else {
		message := "Hello"
		for _, word := range os.Args[1:] {
			message += " " + word
		}
		fmt.Println(message + "\nWelcome to the jungle")
	}
}
