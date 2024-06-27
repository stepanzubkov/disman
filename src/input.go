package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/term"
)

// Get input from console
func getInput(prompt string) string {
    fmt.Print(prompt)
    reader := bufio.NewReader(os.Stdin)
    input, err := reader.ReadString('\n')
    if err != nil {
        log.Fatalln("An error occured while reading input. Please try again", err)
    }

    input = strings.TrimSuffix(input, "\n")
    return input
}

// Get password input from console, hiding user input
func getPasswordInput(prompt string) string {
    fmt.Print(prompt)
    password, err := term.ReadPassword(int(os.Stdin.Fd()))
    if err != nil {
        log.Fatalln("An error occured while reading password input. Please try again", err)
    }
    fmt.Println()
    return string(password)
}
