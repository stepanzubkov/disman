package main

/*
#include <sys/types.h>
#include <pwd.h>
#include <stdlib.h>
*/
import "C"

import (
    "unsafe"
)

type Passwd struct {
    Name   string // user name
    Passwd string // user password
    UID    uint32 // user ID
    GID    uint32 // group ID
    Gecos  string // real name
    Dir    string // home directory
    Shell  string // shell program
}

func cpasswd2go(cpw *C.struct_passwd) *Passwd {
    return &Passwd{
        Name:   C.GoString(cpw.pw_name),
        Passwd: C.GoString(cpw.pw_passwd),
        UID:    uint32(cpw.pw_uid),
        GID:    uint32(cpw.pw_uid),
        Gecos:  C.GoString(cpw.pw_gecos),
        Dir:    C.GoString(cpw.pw_dir),
        Shell:  C.GoString(cpw.pw_shell),
    }
}

// Getpwnam searches the user database for an entry with a matching name.
func Getpwnam(name string) *Passwd {
    cname := C.CString(name)
    defer C.free(unsafe.Pointer(cname))
    cpw := C.getpwnam(cname)
    if cpw == nil {
        return nil
    }
    return cpasswd2go(cpw)
}


