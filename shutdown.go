package main

import (
	"os/exec"
)

func main() {
	if err := exec.Command("cmd", "/C", "shutdown", "/s", "/f", "/t", "0").Run(); err != nil {
		println("Failed to initiate shutdown:", err)
	}
}
