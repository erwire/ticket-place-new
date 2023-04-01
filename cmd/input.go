package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	cmd := exec.Command("./updater_windows_amd64.exe")
	cmdname, _ := os.Executable()

	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "version=0.0.1", fmt.Sprintf("exec_name=%s", cmdname))
	cmd.Env = append(cmd.Env, fmt.Sprintf("pid=%d", os.Getpid()))
	cmd.Env = append(cmd.Env, fmt.Sprintf("repo=%s", "test"))
	cmd.Env = append(cmd.Env, fmt.Sprintf("owner=%s", "jahngeor"))

	var stderr, stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		log.Println("Не пашет")
	}
	fmt.Println(string(stdout.Bytes()), string(stderr.Bytes()))

}
