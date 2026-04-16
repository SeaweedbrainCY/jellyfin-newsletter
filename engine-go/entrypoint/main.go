package main

import (
	"fmt"
	"io/fs"
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

	ChownRecursively(uid, gid, "/app/config")

	DropAndSetNewPrivileges(uid, gid)

	binary, err := exec.LookPath("/app/jellyfin-newsletter")
	if err != nil {
		panic(err)
	}
	// nosemgrep: dangerous-syscall-exec.
	// Expected for an entrypoint.
	_ = syscall.Exec(binary, os.Args, os.Environ())
}

func DropAndSetNewPrivileges(uid, gid int) {
	if err := syscall.Setgid(gid); err != nil {
		panic(err)
	}
	if err := syscall.Setuid(uid); err != nil {
		panic(err)
	}
}

func ChownRecursively(uid, gid int, root string) {
	r, err := os.OpenRoot(root)
	if err != nil {
		panic(err)
	}
	defer r.Close()

	err = fs.WalkDir(r.FS(), ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			panic(err)
		}
		err = r.Lchown(path, uid, gid)
		if err != nil {
			panic(err)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}
