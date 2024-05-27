package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/msteinert/pam/v2"
	"golang.org/x/term"
)

func main() {
    config := parseArgsToConfig()
    var fTTY *os.File
    if (config.Daemon) {
        fTTY = startDaemon()
    }
    fmt.Println("\x1b[01;33m>>> Disman Display Manager <<<\x1b[0m")
    err := errors.New("")
    var t *pam.Transaction
    var username string
    var password string
    for err != nil {
        username = getInput("Username: ")
        password = getPasswordInput("Password: ")
        t, err = checkLogin(username, password)
        if err != nil {
            fmt.Println(err)
        }
    }
    display := ":0"
    vt := "vt7"
    if len(os.Args) == 3 {
        display = os.Args[1]
        vt = os.Args[2]
    }
    initEnv(t, username, display)
    xcmd := startXServer(display, vt, Getpwnam(username))

    sigChan := make(chan os.Signal, 10)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
    go handleKill(xcmd, sigChan)

    cmd := startSession(t, username, display)
    log.Println("Session started")
    cmd.Wait()
    log.Println("Close session")
    stopXServer(xcmd)
    if (config.Daemon) {
        stopDaemon(fTTY)
    }
}

// Handle application kill
func handleKill(Xcmd *exec.Cmd, stopChan chan os.Signal) {
    <- stopChan
    stopXServer(Xcmd)
    log.Fatalln("Exit from an application")
}

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
    return string(password)
}
