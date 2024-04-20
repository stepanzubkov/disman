package main
import (
    "github.com/msteinert/pam/v2"
)

func InitEnv(t *pam.Transaction, pw *Passwd) {
    setEnv(t, "HOME", pw.Dir)
    setEnv(t, "PWD", pw.Dir)
    setEnv(t, "SHELL", pw.Shell)
    setEnv(t, "USER", pw.Name)
    setEnv(t, "LOGNAME", pw.Name)
    setEnv(t, "PATH", "/usr/local/sbin:/usr/local/bin:/usr/bin")
    setEnv(t, "XAUTHORITY", pw.Dir + "/.Xauthority")
}

func setEnv(t *pam.Transaction, name string, value string) {
    name_value := name + "=" + value
    t.PutEnv(name_value)
}
