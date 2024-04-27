package main

import (
	"errors"
	"log"
	"os"
	"os/exec"

	"github.com/msteinert/pam/v2"
)

func checkLogin(login string, password string) (*pam.Transaction, error) {
    t, err := pam.StartFunc("dm", "", conversation(login, password))
    if err != nil {
        return nil, errors.New("PAM start: " + err.Error())
    }
    err = t.Authenticate(0)
    if err != nil {
        return nil, errors.New("PAM auth: " + err.Error())
    }

    err = t.AcctMgmt(0)
    if err != nil {
        return nil, errors.New("PAM acct mgmt: " + err.Error())
    }

    err = t.SetCred(pam.EstablishCred)
    if err != nil {
        return nil, errors.New("PAM set cred: " + err.Error())
    }

    err = t.OpenSession(0)
    if err != nil {
        t.SetCred(pam.DeleteCred)
        return nil, errors.New("PAM open session: " + err.Error())
    }
    return t, nil
}

func startSession(t *pam.Transaction, login string, display string) *exec.Cmd {
    pwd := Getpwnam(login)
    os.Chdir(pwd.Dir)
    log.Println("Start session with user " + login)
    cmd := exec.Command("su", login, "&&", pwd.Shell, "-c", "exec /bin/bash --login .xinitrc")
    cmd.Stdin = os.Stdin
    cmd.Stderr = os.Stderr
    err := cmd.Start()
    if err != nil {
        log.Fatalln(err)
    }
    return cmd
}

func conversation(login string, password string) (func(pam.Style, string) (string, error)) {
    return func (s pam.Style, msg string) (string, error) {
        switch s {
            case pam.PromptEchoOff:
                log.Println(msg)
                return password, nil
            case pam.PromptEchoOn:
                log.Println(msg)
                return login, nil
            case pam.ErrorMsg:
                log.Println("ERROR: ", msg)
                return "", nil
            case pam.TextInfo:
                log.Println(msg)
                return "", nil
            default:
                return "", errors.New("Unrecognized message style")
            }
    }
}
