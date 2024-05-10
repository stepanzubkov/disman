package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"syscall"
	"unsafe"
)

const (
    pathOsReleaseFile = "/etc/os-release"

    osReleasePrettyName = "PRETTY_NAME"
    osReleaseName       = "NAME"

    devPath = "/dev/"

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

// Clears terminal screen
func clearScreen(w io.Writer) {
    if w == nil {
        fmt.Print(strCleanScreen)
    } else {
        w.Write([]byte(strCleanScreen))
    }
}

// Perform switch to defined TTY, if switchTTY is true and tty is greater than 0.
func switchTTY() bool {
    return chvt(7)
}

func openConsole(path string) *os.File {
    for _, flag := range []int{os.O_RDWR, os.O_RDONLY, os.O_WRONLY} {
        if c, err := os.OpenFile(path, flag, 0700); err == nil {
            return c
        }
    }
    return nil
}

// Checks, if used fd is a console
func isConsole(fd uintptr) bool {
    flag := 0
    if _, _, errNo := syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(_KDGKBTYPE), uintptr(unsafe.Pointer(&flag))); errNo == 0 {
        return flag == _KB_101 || flag == _KB_84
    }
    return false
}

func getConsole() *os.File {
    for _, name := range []string{currentTty, currentVc, devConsole} {
        if c := openConsole(name); c != nil {
            if isConsole(c.Fd()) {
                return c
            }
            c.Close()
        }
    }
    return nil
}

func chvt(tty int) bool {
	if c := getConsole(); c != nil {
		defer c.Close()
		if _, _, errNo := syscall.Syscall(syscall.SYS_IOCTL, uintptr(c.Fd()), uintptr(_VT_ACTIVATE), uintptr(tty)); errNo > 0 {
			return false
		}
		if _, _, errNo := syscall.Syscall(syscall.SYS_IOCTL, uintptr(c.Fd()), uintptr(_VT_WAITACTIVE), uintptr(tty)); errNo > 0 {
			return false
		}
	}
	return true
}
