package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/akamensky/argparse"
)

const (
    configFilePath = "/etc/disman.conf"
)

type Config struct {
    Daemon         bool
    Display        string
    Vt             string
    PreCommand     string
    DisplayTitle   bool
    DefaultUser    string
    DefaultSession string
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
        Validate: validateVtArg,
    })
    err := parser.Parse(os.Args)
    if err != nil {
        fmt.Println(parser.Usage(err))
        os.Exit(1)
    }
    // FIXME: Repeating code
    config.Daemon = *daemon
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
        DisplayTitle: true,
        DefaultUser: "",
        DefaultSession: "",
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
            case "DISPLAY":
                // TODO: Logging about invalid value or exit application
                if err := validateDisplayArg([]string{parsedLine.Value}); err == nil {
                    config.Display = parsedLine.Value
                }
            case "VT":
                if err := validateVtArg([]string{parsedLine.Value}); err == nil {
                    config.Vt = parsedLine.Value
                }
            case "PRE_COMMAND":
                config.PreCommand = parsedLine.Value
            case "DISPLAY_TITLE":
                config.DisplayTitle = parseBool(parsedLine.Value, true)
            case "DEFAULT_USER":
                config.DefaultUser = parsedLine.Value
            case "DEFAULT_SESSION":
                config.DefaultSession = parsedLine.Value
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


func validateVtArg(args []string) error {
    commonError := errors.New("Invalid virtual terminal name!")
    vt := args[0]
    if vt[:2] != "vt" {
        return commonError
    }
    if n, err := strconv.Atoi(vt[2:]); err != nil || n < 1 {
        return commonError
    }
    return nil
}


// Parses bool value, returns defaultValue as default (If value is not a valid boolean)
func parseBool(value string, defaultValue bool) bool {
    res, err := strconv.ParseBool(value)
    if err != nil {
        return defaultValue
    }
    return res
}


// Runs PRE_COMMAND defined in config file
func runPreCommand(config *Config) {
    cmd := exec.Command("/bin/bash", "-c", config.PreCommand) 
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    err := cmd.Run()
    if err != nil {
        fmt.Printf("Error while executing PRE_COMMAND: %v", err)
    }
}
