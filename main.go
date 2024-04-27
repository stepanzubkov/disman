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
	"time"

	"github.com/msteinert/pam/v2"
)

const testing = false

func main() {
    fmt.Println("--- Small Display Manager ---")
    err := errors.New("")
    var t *pam.Transaction
    var username string
    var password string
    for err != nil {
        username = getInput("Username: ")
        password = getInput("Password: ")
        t, err = checkLogin(username, password)
        if err != nil {
            fmt.Println(err)
        }
    }
    display := ":0"
    vt := "vt1"
    if len(os.Args) == 3 {
        display = os.Args[1]
        vt = os.Args[2]
    }
    log.Println("Before X server started")
    xcmd := startXServer(display, vt)
    log.Println("X server started!")

    sigChan := make(chan os.Signal, 10)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
    go handleKill(xcmd, sigChan)

    initEnv(t, username, display)
    cmd := startSession(t, username, display)
    log.Println("Session started")
    cmd.Wait()
    log.Println("Close session")
    stopXServer(xcmd)
}


func startXServer(display string, vt string) *exec.Cmd {
    cmd := exec.Command("/usr/bin/X", display, vt)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    err := cmd.Start()
    if err != nil {
        log.Fatalln(err)
    }
    time.Sleep(3 * time.Second)
    if err != nil {
        log.Fatalln("Xorg: ", err)
    }
    return cmd
}

func stopXServer(Xcmd *exec.Cmd) {
    Xcmd.Process.Kill()
}

func handleKill(Xcmd *exec.Cmd, stopChan chan os.Signal) {
    <- stopChan
    stopXServer(Xcmd)
    log.Fatalln("Exit from an application")
}

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
