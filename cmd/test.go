package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	pwd, _ := os.Getwd()

	fmt.Println(pwd)
	cmdPath := fmt.Sprintf("G:\\Проекты\\Работа\\Рабочие проекты\\Freelance\\Григорий\\ticket-place-new\\build\\ticket-place_windows_amd64.exe")
	cmd := exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", cmdPath)

	fmt.Println(cmd.Run())
}
