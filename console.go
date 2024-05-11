package main

import (
	"fmt"
	"io"
	"os"
	"syscall"
	"unsafe"
)

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

// Opens console file
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

// Opens console and checks if it is console
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

// Change virtual terminal
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
