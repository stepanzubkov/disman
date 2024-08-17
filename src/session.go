package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
    lastSessionDir = "/.cache/disman"
    lastSessionFile = "lastsession"
    lastSessionPath = lastSessionDir + "/" + lastSessionFile
)

// Returns command needed for starting X session
func getSessionEntry(user *User, config *Config) *DesktopEntry {
    sessions := getSessions()

    defaultSession := getDefaultSession(sessions, config) 
    if defaultSession != nil {
        return defaultSession
    }

    inputLabel := "Choose session "
    for index, session := range sessions {
        inputLabel = inputLabel + fmt.Sprintf("[%d] %s", index + 1, session.Name)
        if index != len(sessions) - 1 {
            inputLabel = inputLabel + ", "
        }
    }
    lastSession := getLastSession(user, sessions)
    if lastSession != -1 {
        inputLabel = inputLabel + fmt.Sprintf(" (%v): ", lastSession + 1)
    } else {
        inputLabel = inputLabel + ": "
    }
    var sessionNumber int
    var err error
    for {
        sessionString := getInput(inputLabel)
        if sessionString == "" && lastSession != -1 {
            return sessions[lastSession]
        }
        sessionNumber, err = strconv.Atoi(sessionString) 
        if err != nil {
            fmt.Println("Your input is not integer!")
            continue
        }
        if sessionNumber < 1 || sessionNumber > len(sessions) {
            fmt.Printf("Choose number from 1 to %d\n", len(sessions))
            continue
        }
        break
    }
    choosedSession := sessions[sessionNumber-1]
    return choosedSession
}


func getSessions() []*DesktopEntry {
    files, err := os.ReadDir("/usr/share/xsessions")
    if err != nil {
        log.Fatalf("Can't read dir /usr/share/xsessions! %s\n", err)
    }

    var desktopEntries []*DesktopEntry
    for _, file := range files {
        if !file.IsDir() && strings.HasSuffix(file.Name(), ".desktop") {
            desktopEntries = append(desktopEntries, parseDesktopEntry("/usr/share/xsessions/" + file.Name()))
        }
    }
    return desktopEntries
}


// Searches default session in all sessions by its Name field. If default session is not found returns nil.
func getDefaultSession(sessions []*DesktopEntry, config *Config) *DesktopEntry {
    if config.DefaultSession == "" {
        return nil
    }
    for _, session := range sessions {
        if session.Name == config.DefaultSession {
            return session
        }
    }
    return nil
}


// Gets last logged in session index in session array for user
func getLastSession(user *User, sessions []*DesktopEntry) int {
    lastSessionPathForUser := user.Dir + lastSessionPath
    _, err := os.Stat(lastSessionPathForUser)
    if err != nil {
        return -1
    }
    fileContent, err := os.ReadFile(lastSessionPathForUser)
    lastSessionExec := strings.TrimSpace(string(fileContent))
    for index, session := range sessions {
        if session.Exec == lastSessionExec {
            return index
        }
    }
    return -1
}

// Writes session as last logged in session
func writeLastSession(session *DesktopEntry, user *User) {
    lastSessionPathForUser := user.Dir + lastSessionPath
    lastSessionDirForUser := user.Dir + lastSessionDir
    err := os.MkdirAll(lastSessionDirForUser, 0755)
    if err != nil {
        log.Fatalf("Unable to create last session directory! %v\n", err)
    }
    err = os.WriteFile(lastSessionPathForUser, []byte(session.Exec), 0644)
    if err != nil {
        log.Fatalf("Unable to create last session file! %v\n", err)
    }
}
