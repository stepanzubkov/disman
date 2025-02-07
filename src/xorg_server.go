package main

import (
	"log"
	"os"
	"os/exec"
	"time"
)

// Starts Xorg server
func startXServer(config *Config, user *User) *exec.Cmd {
    cmd := exec.Command("/bin/bash", "-c", "/usr/bin/X " + config.Display + " " + config.Vt)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    os.Setenv("DISPLAY", config.Display)
    os.Setenv("XAUTHORITY", user.Dir + "/.Xauthority")
    cmd.Env = os.Environ()
    err := cmd.Start()
    if err != nil {
        log.Fatalln(err)
    }
    // TODO: Wait until Xorg started
    time.Sleep(3 * time.Second)
    return cmd
}

// Stops Xorg server by sending interrupt signal
func stopXServer(Xcmd *exec.Cmd) {
    Xcmd.Process.Signal(os.Interrupt)
    log.Println("Stop Xorg")
    Xcmd.Wait()
}
