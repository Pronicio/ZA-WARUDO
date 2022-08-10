package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	println("Starting...")

	dir, err := os.MkdirTemp("", "za-warudo")
	if err != nil {
		log.Fatal(err)
	}

	if err := downloadFile(dir, "https://za-warudo.vercel.app/video.mp4", "video", "mp4"); err != nil {
		println("Failed :", err)
	}
	
	if err := downloadFile(dir, "https://za-warudo.vercel.app/ffplay.zip", "video", "mp4"); err != nil {
		println("Failed :", err)
	}

	/*
		if err := exec.Command("cmd", "/C", "ffplay.exe", "video.mp4").Run(); err != nil {
			println("Failed :", err)
		}
	*/

	println("End")
}

func downloadFile(filepath string, url string, filename string, filetype string) (err error) {

	// Create the file
	out, err := os.CreateTemp(filepath, filename)
	if err != nil {
		log.Fatal(err)
	}

	println(out.Name())

	defer out.Close()

	// Get the data
	res, err := http.Get(url)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	// Check server response
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", res.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, res.Body)
	if err != nil {
		return err
	}

	return nil
}
