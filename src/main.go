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
    lastUser := getLastUser()
    for err != nil {
        if lastUser != nil {
            username = getInput("Username (" + lastUser.Name + "): ")
            if username == "" {
                username = lastUser.Name
            }
        } else {
            username = getInput("Username: ")
        }
        password = getPasswordInput("Password: ")
        t, err = checkLogin(username, password)
        if err != nil {
            fmt.Println(err)
        }
    }
    user := getUser(username)
    user.writeLastUser()

    sessionEntry := getSessionEntry(user)
    writeLastSession(sessionEntry, user)

    initEnv(t, user, config, sessionEntry)
    xcmd := startXServer(config, user)

    sigChan := make(chan os.Signal, 10)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
    go handleKill(xcmd, sigChan)

    cmd := startSession(t, username, "exec " + sessionEntry.Exec)
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
