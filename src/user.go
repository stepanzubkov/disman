package main

/*
#include <sys/types.h>
#include <pwd.h>
#include <stdlib.h>
*/
import "C"

import (
	"log"
	"os"
	"os/user"
	"strconv"
	"strings"
	"unsafe"
)

const (
    lastUsernameDir = "/var/cache/disman"
    lastUsernameFile = "lastuser"
    lastUsernamePath = lastUsernameDir + "/" + lastUsernameFile
)

type User struct {
    // Passwd fields
    Name   string   // user name
    Passwd string   // user password
    UID    uint32   // user ID
    GID    uint32   // group ID
    Gecos  string   // real name
    Dir    string   // home directory
    Shell  string   // shell program

    // Other fields
    Gids   []uint32 // user's groups ids
}


// Returns last logged in user
func getLastUser() *User {
    _, err := os.Stat(lastUsernamePath)
    if err != nil {
        return nil
    }
    fileContent, err := os.ReadFile(lastUsernamePath)
    username := strings.TrimSpace(string(fileContent))
    user := getUser(username)
    return user
}

// Writes user as last logged in user
func (user *User) writeLastUser() {
    err := os.MkdirAll(lastUsernameDir, 0755)
    if err != nil {
        log.Fatalf("Unable to create last user directory! %v\n", err)
    }
    err = os.WriteFile(lastUsernamePath, []byte(user.Name), 0644)
    if err != nil {
        log.Fatalf("Unable to create last user file! %v\n", err)
    }
}


// Returns full user structure
func getUser(username string) *User {
    user := getUserFromPasswd(username)
    if user == nil {
        return nil
    }
    user.Gids = getUserGids(username)
    return user
}


// Searches for username in passwd file and returns User structure
func getUserFromPasswd(username string) *User {
    cname := C.CString(username)
    defer C.free(unsafe.Pointer(cname))
    cpw := C.getpwnam(cname)
    if cpw == nil {
        return nil
    }
    return cpasswd2go(cpw)
}

// Returns user's groups ids
func getUserGids(username string) []uint32 {
    var gids []uint32
    user, _ := user.Lookup(username)
    if strGids, err := user.GroupIds(); err == nil {
        for _, val := range strGids {
            value, _ := strconv.Atoi(val)
            gids = append(gids, uint32(value))
        }
    }
    return gids
}


// Converts C passwd structure to Go User structure
func cpasswd2go(cpw *C.struct_passwd) *User {
    return &User{
        Name:   C.GoString(cpw.pw_name),
        Passwd: C.GoString(cpw.pw_passwd),
        UID:    uint32(cpw.pw_uid),
        GID:    uint32(cpw.pw_uid),
        Gecos:  C.GoString(cpw.pw_gecos),
        Dir:    C.GoString(cpw.pw_dir),
        Shell:  C.GoString(cpw.pw_shell),
    }
}
