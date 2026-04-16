package main

import (
	"fmt"
	"os"
	"os/exec"
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

	if err := syscall.Setgid(gid); err != nil {
		panic(err)
	}
	if err := syscall.Setuid(uid); err != nil {
		panic(err)
	}

	binary, _ := exec.LookPath("/app/jellyfin-newsletter")
	syscall.Exec(binary, os.Args, os.Environ())
}
