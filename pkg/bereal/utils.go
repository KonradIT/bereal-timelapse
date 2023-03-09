package bereal

import (
	"fmt"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/dustin/go-humanize"
	"golang.org/x/image/webp"
)

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
	// Clear the line by using a character return to go back to the start and remove
	// the remaining characters by filling it with spaces
	fmt.Printf("\r%s", strings.Repeat(" ", 35))

	// Return again and print current status of download
	// We use the humanize package to print the bytes in a meaningful way (e.g. 10 MB)
	fmt.Printf("\rDownloading... %s complete", humanize.Bytes(wc.Total))
}

func DownloadFile(localpath string, url string) error {
	// Create the file, but give it a tmp file extension, this means we won't overwrite a
	// file until it's downloaded, but we'll remove the tmp extension once downloaded.
	out, err := os.Create(localpath + ".tmp")
	if err != nil {
		return err
	}

	// Get the data
	resp, err := http.Get(url) //nolint
	if err != nil {
		out.Close()
		return err
	}
	defer resp.Body.Close()

	// Create our progress reporter and pass it to be used alongside our writer
	counter := &WriteCounter{}
	if _, err = io.Copy(out, io.TeeReader(resp.Body, counter)); err != nil {
		out.Close()
		return err
	}

	// The progress use the same line so print a new line once it's finished downloading
	fmt.Print("\n")

	// Close the file without defer so it can happen before Rename()
	out.Close()

	extension := filepath.Ext(url)
	switch extension {
	case ".webp":
		f0, err := os.Open(localpath + ".tmp")
		if err != nil {
			return err
		}

		jpgImg, err := os.Create(localpath + ".tmp.c")
		if err != nil {
			return err
		}
		img0, err := webp.Decode(f0)
		if err != nil {
			return err
		}
		err = jpeg.Encode(jpgImg, img0, &jpeg.Options{
			Quality: jpeg.DefaultQuality,
		})
		if err != nil {
			return err
		}
		if err := f0.Close(); err != nil {
			return err
		}
		if err := jpgImg.Close(); err != nil {
			return err
		}
		if err := os.Rename(localpath+".tmp.c", localpath+".tmp"); err != nil {
			return err
		}
	}

	if err = os.Rename(localpath+".tmp", localpath); err != nil {
		return err
	}
	return nil
}
