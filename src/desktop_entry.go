package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

type DesktopEntryLine struct {
    Name  string;
    // Value may be number or string
    Value string;
};

type DesktopEntry struct {
    // Not all keys from specification, just all keys needed for xsession

    Type         string;
    Exec         string;
    TryExec      string;
    DesktopNames string;
    Name         string;
    Comment      string;
}

// Parse file with .desktop extension (desktop entry)
func parseDesktopEntry(path string) *DesktopEntry {
    file, err := os.Open(path)
    if err != nil {
        log.Fatalf("Unable to open desktop entry %s. %s\n", path, err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    desktopEntry := &DesktopEntry{}
    for scanner.Scan() {
        parsedLine := parseLine(scanner.Text())
        if parsedLine == nil {
            continue
        }
        switch parsedLine.Name {
            case "Type":
                desktopEntry.Type = parsedLine.Value
            case "Exec":
                desktopEntry.Exec = parsedLine.Value
            case "TryExec":
                desktopEntry.TryExec = parsedLine.Value
            case "DesktopNames":
                desktopEntry.DesktopNames = parsedLine.Value
            case "Name":
                desktopEntry.Name = parsedLine.Value
            case "Comment":
                desktopEntry.Comment = parsedLine.Value
        }
    }
    if err = scanner.Err(); err != nil {
        log.Fatalln(err)
    }

    return desktopEntry
}

// Parser a line of desktop entry
func parseLine(line string) *DesktopEntryLine {
    line = strings.Trim(line, " ")

    // Just skip empty and comment lines
    // Also skip section headers
    // TODO: Do not skip section headers
    if line == "" || line[0] == '#' || line[0] == '[' && line[len(line)-1] == ']' {
        return nil
    }

    // In FreeDesktop.org specification, '=' char not allowed in name or value,
    // so line can contain only one '=' char as delimeter
    lineSplit := strings.Split(line, "=")
    if len(lineSplit) != 2 {
        log.Fatalf("Invalid line in desktop entry: \"%s\".\n", line)
    }
    name := strings.TrimRight(lineSplit[0], " ")
    value := strings.TrimLeft(lineSplit[1], " ")
    return &DesktopEntryLine{Name: name, Value: value}
}
