package main

import (
	"fmt"
	"os"

	"github.com/akamensky/argparse"
)

type Config struct {
    Daemon bool
}

// Parse os.Args to Config struct
func parseArgsToConfig() *Config {
    parser := argparse.NewParser("disman", "CLI Display Manager")
    daemon := parser.Flag("d", "daemon", &argparse.Options{
        Required: false,
        Help: "Run as daemon",
        Default: false,
    })
    err := parser.Parse(os.Args)
    if err != nil {
        fmt.Println(parser.Usage(err))
        os.Exit(1)
    }
    return &Config{
        Daemon: *daemon,
    }
}
