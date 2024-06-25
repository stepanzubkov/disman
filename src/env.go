package main

import (
	"log"
	"os"
	"strconv"

	"github.com/msteinert/pam/v2"
)

// Initializes environment for X session
func initEnv(t *pam.Transaction, login string, display string) {
    passwd := Getpwnam(login)
    setEnv(t, "HOME", passwd.Dir)
    setEnv(t, "PWD", passwd.Dir)
    setEnv(t, "SHELL", passwd.Shell)
    setEnv(t, "USER", passwd.Name)
    setEnv(t, "LOGNAME", passwd.Name)
    setEnv(t, "PATH", "/usr/local/sbin:/usr/local/bin:/usr/bin")
    setEnv(t, "XAUTHORITY", passwd.Dir + "/.Xauthority")
    setEnv(t, "DISPLAY", display)
    xdg_runtime_dir := "/run/user/" + strconv.FormatUint(uint64(passwd.UID), 10)
    setEnv(t, "XDG_RUNTIME_DIR", xdg_runtime_dir)

    createXdgRuntimeDir(xdg_runtime_dir, passwd)
}

// Create XDG_RUNTIME_DIR if needed
func createXdgRuntimeDir(dir string, passwd *Passwd) {
    err := os.MkdirAll(dir, 0700)
    if err != nil {
        log.Fatalf("Unable to create XDG_RUNTIME_DIR! %s\n", err)
    }

    err = os.Chown(dir, int(passwd.UID), int(passwd.GID))
    if err != nil {
        log.Fatalf("Unable to change owner of XDG_RUNTIME_DIR! %s\n", err)
    }
}

// Sets environment variable for PAM transaction
func setEnv(t *pam.Transaction, name string, value string) {
    name_value := name + "=" + value
    t.PutEnv(name_value)
}

