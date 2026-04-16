package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
)

func main() {
	uidStr := "1001"
	gidStr := "1001"
	if uid := os.Getenv("USER_UID"); uid != "" {
		uidStr = uid
	} else {
		fmt.Println("USER_UID environment variable not set. Using default value: " + uidStr)
	}

	if gid := os.Getenv("USER_GID"); gid != "" {
		gidStr = gid
	} else {
		fmt.Println("USER_GID environment variable not set. Using default value: " + gidStr)
	}
	uid, err := strconv.Atoi(uidStr)
	if err != nil {
		panic(err)
	}

	gid, err := strconv.Atoi(gidStr)
	if err != nil {
		panic(err)
	}

	ChownRecursively(uid, gid, "/app/config")

	DropAndSetNewPrivileges(uid, gid)

	binary, _ := exec.LookPath("/app/jellyfin-newsletter")
	syscall.Exec(binary, os.Args, os.Environ())
}

func DropAndSetNewPrivileges(uid, gid int) {
	if err := syscall.Setgid(gid); err != nil {
		panic(err)
	}
	if err := syscall.Setuid(uid); err != nil {
		panic(err)
	}
}

// Source - https://stackoverflow.com/a/73864967
// Posted by h0ch5tr4355
// Retrieved 2026-04-16, License - CC BY-SA 4.0
func ChownRecursively(uid, gid int, root string) {
	err := filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				panic(err)
			}
			err = os.Chown(path, uid, gid)
			if err != nil {
				panic(err)
			}
			return nil
		})
	if err != nil {
		panic(err)
	}
}
