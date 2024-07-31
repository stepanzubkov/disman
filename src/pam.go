package main

import (
	"errors"
	"log"
    "fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/msteinert/pam/v2"
)

// Asks user for username/password pair and returns pam transaction and username
func getLoginCredentialsFromUser() (t *pam.Transaction, username string) {
    err := errors.New("")
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
    return t, username
}

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
    user := getUser(login)
    os.Chdir(user.Dir)
    log.Println("Start session with user " + login)
    cmd := exec.Command(user.Shell, "-c", command)
    prepareSessionCmd(cmd, user, t)

    err := cmd.Start()
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


// Sets neccessary flags for running session command
func prepareSessionCmd(cmd *exec.Cmd, user *User, t *pam.Transaction) {
    cmd.Stdin = os.Stdin
    cmd.Stderr = os.Stderr
    cmd.Stdout = os.Stdout

    // Flags to run as user
    cmd.SysProcAttr = &syscall.SysProcAttr{}
    cmd.SysProcAttr.Credential = &syscall.Credential{Uid: user.UID, Gid: user.GID, Groups: user.Gids}

    copyPamEnvToCmdEnv(cmd, t)
}

// Copies PAM transaction env to cmd env
func copyPamEnvToCmdEnv(cmd *exec.Cmd, t *pam.Transaction) {
    var envList []string
    envMap, err := t.GetEnvList()
    if err != nil {
        log.Fatalf("Can't get env list of pam transaction! %s\n", err)
    }
    for key, value := range envMap {
        envList = append(envList, key+"="+value)
    }
    cmd.Env = envList
}
