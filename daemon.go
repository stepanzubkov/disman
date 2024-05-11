package main

import (
	"log"
	"os"
)

const (
    _KDGKBTYPE     = 0x4B33
    _VT_ACTIVATE   = 0x5606
    _VT_WAITACTIVE = 0x5607

    _KB_101 = 0x02
    _KB_84  = 0x01

    currentTty = "/dev/tty"
    currentVc = "/dev/tty0"
    devConsole = "/dev/console"

    strCleanScreen = "\x1b[H\x1b[2J"
)

// Opens other tty and switch to it
func startDaemon() *os.File {
    fTTY, err := os.OpenFile("/dev/tty7", os.O_RDWR, 0700)
    if err != nil {
        log.Fatalln(err)
    }

    clearScreen(fTTY)

    os.Stdout = fTTY
    os.Stderr = fTTY
    os.Stdin = fTTY

    clearScreen(fTTY)

    switchTTY()

    return fTTY
}

// Stops daemon mode and closes opened TTY, if allowed
func stopDaemon(fTTY *os.File) {
    clearScreen(fTTY)

    if fTTY != nil {
        fTTY.Close()
    }
}
