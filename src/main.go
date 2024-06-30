package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/msteinert/pam/v2"
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

    sessionCommand := getSessionCommand()

    initEnv(t, username, config)
    xcmd := startXServer(config, getUser(username))

    sigChan := make(chan os.Signal, 10)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
    go handleKill(xcmd, sigChan)

    cmd := startSession(t, username, sessionCommand)
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
