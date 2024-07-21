package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/akamensky/argparse"
)

const (
    configFilePath = "/etc/disman.conf"
)

type Config struct {
    Daemon  bool
    Display string
    Vt      string
}


func parseConfig() *Config {
    baseConfig := parseConfigFileToConfig()
    additionalConfig := extendConfigWithArgs(baseConfig)
    return additionalConfig
}

// Extend config with os.Args
func extendConfigWithArgs(config *Config) *Config {
    parser := argparse.NewParser("disman", "CLI Display Manager")
    daemon := parser.Flag("d", "daemon", &argparse.Options{
        Required: false,
        Help: "Run as daemon",
    })
    display := parser.String("D", "display", &argparse.Options{
        Required: false,
        Help: "X display name",
        Validate: validateDisplayArg,
    })
    vt := parser.String("v", "vt", &argparse.Options{
        Required: false,
        Help: "Virtual terminal number (in form 'vtX')",
    })
    err := parser.Parse(os.Args)
    if err != nil {
        fmt.Println(parser.Usage(err))
        os.Exit(1)
    }
    // FIXME: Repeating code
    config.Daemon = config.Daemon || *daemon
    if *display != "" {
        config.Display = *display
    }
    if *vt != "" {
        config.Vt = *vt
    }
    return config
}


func parseConfigFileToConfig() *Config {
    config := &Config{
        Display: ":0",
        Vt: "vt7",
        Daemon: false,
    }
    file, err := os.Open(configFilePath)
    if err != nil {
        return config
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    // Initialize config with default values
    for scanner.Scan() {
        parsedLine := parseLine(scanner.Text())
        if parsedLine == nil {
            continue
        }
        switch parsedLine.Name {
            // TODO: DAEMON option
            case "DISPLAY":
                // TODO: Validation
                config.Display = parsedLine.Value
            case "VT":
                // TODO: Validation
                config.Vt = parsedLine.Value
        }
    }
    if err = scanner.Err(); err != nil {
        log.Fatalln(err)
    }

    return config
}


// Validates X display name, in format '[host]:<display>.[screen]'
func validateDisplayArg(args []string) error {
    commonError := errors.New("Invalid display name!")
    displayFull := args[0]
    displaySplit := strings.Split(displayFull, ":")
    if len(displaySplit) != 2 {
        return commonError
    }
    displayWithoutHost := displaySplit[1] 
    displayAndScreen := strings.Split(displayWithoutHost, ".")
    if len(displayAndScreen) > 2 {
        return commonError
    }
    for _, s := range displayAndScreen {
        if n, err := strconv.Atoi(s); err != nil || n < 0 {
            return commonError
        }
    }
    return nil
}
