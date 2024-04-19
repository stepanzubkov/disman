package main

import (
    "log"
    "errors"

    "github.com/msteinert/pam"
)

func Login(login string, password string) {
    t, err := pam.StartFunc("", "", conversation(login, password))
    if err != nil {
        log.Fatalln("PAM start: ", err)
    }
    err = t.Authenticate(0)
    if err != nil {
        log.Fatalln("PAM authenticate: ", err)
    }

    err = t.AcctMgmt(0)
    if err != nil {
        log.Fatalln("PAM account management: ", err)
    }

    err = t.SetCred(pam.EstablishCred)
    if err != nil {
        log.Fatalln("PAM set cred: ", err)
    }

    err = t.OpenSession(0)
    if err != nil {
        t.SetCred(pam.DeleteCred)
        log.Fatalln("PAM open session: ", err)
    }
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
