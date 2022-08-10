package main

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	println("Starting...")

	dir, err := os.MkdirTemp("", "za-warudo")
	if err != nil {
		println(err)
	}

	videoPath, err := downloadFile(dir, "https://za-warudo.vercel.app/video.mp4", "video", "mp4")
	if err != nil {
		println("Failed to download the video :", err)
	} else {
		println("videoPath", videoPath)
	}

	zipPath, err := downloadFile(dir, "https://za-warudo.vercel.app/ffplay.zip", "test", "zip")
	if err != nil {
		println("Failed to download the video :", err)
	} else {
		println("zipPath", zipPath)
		err := extractZip(zipPath, dir)
		if err != nil {
			return
		}
	}

	if err := exec.Command("cmd", "/C", dir+"\\"+"ffplay.exe", videoPath, "-autoexit").Run(); err != nil {
		println("Failed :", err)
	}

	println("End")
}

func downloadFile(filepath string, url string, filename string, filetype string) (string, error) {

	// Create the file
	out, err := os.CreateTemp(filepath, filename)
	if err != nil {
		println(err)
	}

	defer out.Close()

	// Get the data
	res, err := http.Get(url)
	if err != nil {
		return "null", err
	}

	defer res.Body.Close()

	// Check server response
	if res.StatusCode != http.StatusOK {
		return "null", fmt.Errorf("bad status: %s", res.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, res.Body)
	if err != nil {
		return "null", err
	}

	return out.Name(), nil
}

func extractZip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(dest, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", path)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}
