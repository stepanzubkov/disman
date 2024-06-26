package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// Returns command needed for starting X session
func getSessionCommand() string {
    sessions := getSessions()
    inputLabel := "Choose session ("
    for index, session := range sessions {
        inputLabel = inputLabel + fmt.Sprintf("[%d] %s", index+1, session.Name)
        if index == len(sessions) - 1 {
            inputLabel = inputLabel + "): "
        } else {
            inputLabel = inputLabel + ", "
        }
    }
    var sessionNumber int
    var err error
    for {
        sessionString := getInput(inputLabel)
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
    return "exec " + choosedSession.Exec
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
