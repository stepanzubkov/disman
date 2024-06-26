package main

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"syscall"

	"github.com/msteinert/pam/v2"
)

// Checks login/password pair with PAM
func checkLogin(login string, password string) (*pam.Transaction, error) {
    t, err := pam.StartFunc("disman", "", conversation(login, password))
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

// Starts X session
func startSession(t *pam.Transaction, login string, command string) *exec.Cmd {
    pwd := Getpwnam(login)
    os.Chdir(pwd.Dir)
    log.Println("Start session with user " + login)
    cmd := exec.Command(pwd.Shell, "-c", command)
    cmd.Stdin = os.Stdin
    cmd.Stderr = os.Stderr
    cmd.Stdout = os.Stdout

    var envList []string
    envMap, err := t.GetEnvList()
    if err != nil {
        log.Fatalf("Can't get env list of pam transaction! %s\n", err)
    }
    for key, value := range envMap {
        envList = append(envList, key+"="+value)
    }
    cmd.Env = envList

    var gids []uint32
    user, _ := user.Lookup(login)
    if strGids, err := user.GroupIds(); err == nil {
        for _, val := range strGids {
            value, _ := strconv.Atoi(val)
            gids = append(gids, uint32(value))
        }
    }
    cmd.SysProcAttr = &syscall.SysProcAttr{}
    cmd.SysProcAttr.Credential = &syscall.Credential{Uid: pwd.UID, Gid: pwd.GID, Groups: gids}

    err = cmd.Start()
    if err != nil {
        log.Fatalln(err)
    }
    return cmd
}

// Recieves messages from PAM
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
