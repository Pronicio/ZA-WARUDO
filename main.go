package main

import (
	"os/exec"
)

func main() {
	println("Starting...")

	if err := exec.Command("cmd", "/C", "ffplay.exe", "video.mp4").Run(); err != nil {
		println("Failed :", err)
	}
	
	println("End")
}
