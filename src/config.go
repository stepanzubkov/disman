package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/akamensky/argparse"
)

type Config struct {
    Daemon  bool
    Display string
    Vt      string
}

// Parse os.Args to Config struct
func parseArgsToConfig() *Config {
    parser := argparse.NewParser("disman", "CLI Display Manager")
    daemon := parser.Flag("d", "daemon", &argparse.Options{
        Required: false,
        Help: "Run as daemon",
        Default: false,
    })
    display := parser.String("D", "display", &argparse.Options{
        Required: false,
        Help: "X display name",
        Default: ":0",
        Validate: validateDisplayArg,
    })
    vt := parser.String("v", "vt", &argparse.Options{
        Required: false,
        Help: "Virtual terminal number (in form 'vtX')",
        Default: "vt7",
    })
    err := parser.Parse(os.Args)
    if err != nil {
        fmt.Println(parser.Usage(err))
        os.Exit(1)
    }
    return &Config{
        Daemon: *daemon,
        Display: *display,
        Vt: *vt,
    }
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
