package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func main() {
    config := parseConfig()
    var fTTY *os.File
    if (config.Daemon) {
        fTTY = startDaemon()
    }
    runPreCommand(config)
    if config.DisplayTitle {
        fmt.Println("\x1b[01;33m>>> Disman Display Manager <<<\x1b[0m")
    }
    t, username := getLoginCredentialsFromUser()
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
