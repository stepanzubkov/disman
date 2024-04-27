package main
import (
    "strconv"
    "github.com/msteinert/pam/v2"
)

func initEnv(t *pam.Transaction, login string, display string) {
    pw := Getpwnam(login)
    setEnv(t, "HOME", pw.Dir)
    setEnv(t, "PWD", pw.Dir)
    setEnv(t, "SHELL", pw.Shell)
    setEnv(t, "USER", pw.Name)
    setEnv(t, "LOGNAME", pw.Name)
    setEnv(t, "PATH", "/usr/local/sbin:/usr/local/bin:/usr/bin")
    setEnv(t, "XAUTHORITY", pw.Dir + "/.Xauthority")
    setEnv(t, "DISPLAY", display)
    setEnv(t, "XDG_RUNTIME_DIR", "/run/user/" + strconv.FormatUint(uint64(pw.UID), 10))
}

func setEnv(t *pam.Transaction, name string, value string) {
    name_value := name + "=" + value
    t.PutEnv(name_value)
}
