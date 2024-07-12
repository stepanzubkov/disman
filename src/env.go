package main

import (
	"log"
	"os"
	"strconv"

	"github.com/msteinert/pam/v2"
)

// Initializes environment for X session
func initEnv(t *pam.Transaction, login string, config *Config, desktopEntry *DesktopEntry) {
    user := getUser(login)
    setEnv(t, "HOME", user.Dir)
    setEnv(t, "PWD", user.Dir)
    setEnv(t, "SHELL", user.Shell)
    setEnv(t, "USER", user.Name)
    setEnv(t, "LOGNAME", user.Name)
    setEnv(t, "PATH", "/usr/local/sbin:/usr/local/bin:/usr/bin")
    setEnv(t, "XAUTHORITY", user.Dir + "/.Xauthority")
    setEnv(t, "DISPLAY", config.Display)

    setEnvIfEmpty(t, "XDG_CONFIG_HOME", user.Dir + "/.config")
    setEnvIfEmpty(t, "XDG_SEAT", "seat0")
    setEnv(t, "XDG_SESSION_CLASS", "user")
    xdg_runtime_dir := "/run/user/" + strconv.FormatUint(uint64(user.UID), 10)
    setEnvIfEmpty(t, "XDG_RUNTIME_DIR", xdg_runtime_dir)
    createXdgRuntimeDir(t.GetEnv("XDG_RUNTIME_DIR"), user)

    // It is deprecated env variable
    // See https://superuser.com/questions/1074068/what-is-the-difference-between-desktop-session-xdg-session-desktop-and-xdg-cur
    setEnv(t, "DESKTOP_SESSION", desktopEntry.Name)
    setEnv(t, "XDG_SESSION_DESKTOP", desktopEntry.getDesktopName())
    if desktopEntry.DesktopNames != "" {
        setEnv(t, "XDG_CURRENT_DESKTOP", desktopEntry.getDesktopName())
    }


}

// Create XDG_RUNTIME_DIR if needed
func createXdgRuntimeDir(dir string, user *User) {
    err := os.MkdirAll(dir, 0700)
    if err != nil {
        log.Fatalf("Unable to create XDG_RUNTIME_DIR! %s\n", err)
    }

    err = os.Chown(dir, int(user.UID), int(user.GID))
    if err != nil {
        log.Fatalf("Unable to change owner of XDG_RUNTIME_DIR! %s\n", err)
    }
}

// Sets environment variable for PAM transaction
func setEnv(t *pam.Transaction, name string, value string) {
    name_value := name + "=" + value
    t.PutEnv(name_value)
}

func setEnvIfEmpty(t *pam.Transaction, name string, value string) {
    if t.GetEnv(name) == "" {
        setEnv(t, name, value)
    }
}

