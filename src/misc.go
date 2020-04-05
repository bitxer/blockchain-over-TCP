package main

import (
	"fmt"
	"strings"
)

func printSuccess(message ...string) {
	fmt.Printf("\033[1;32m\033[1;1m[+]\033[1;0m %s\n", strings.Join(message, " "))
}

func printInfo(message ...string) {
	fmt.Printf("\033[1;33m\033[1;1m[*]\033[1;0m %s\n", strings.Join(message, " "))
}

func printError(message ...string) {
	fmt.Printf("\033[1;31m\033[1;1m[-]\033[1;0m %s\n", strings.Join(message, " "))
}

func printPrompt(message ...string) {
	fmt.Printf("\033[1;33m\033[1;1m[*]\033[1;0m %s ", strings.Join(message, " "))
}
