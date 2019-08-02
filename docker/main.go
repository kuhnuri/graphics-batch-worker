package main

import (
	"fmt"
	kuhnuri "github.com/kuhnuri/go-worker"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
)

type Args struct {
	src *url.URL
	dst *url.URL
	tmp string
	out string
}

func readArgs() *Args {
	input := os.Getenv("input")
	if input == "" {
		log.Fatalf("Input environment variable not set")
	}
	output := os.Getenv("output")
	if output == "" {
		log.Fatalf("Output environment variable not set")
	}
	src, err := url.Parse(input)
	if err != nil {
		log.Fatalf("Failed to parse input argument %s: %v", input, err)
	}
	dst, err := url.Parse(output)
	if err != nil {
		log.Fatalf("Failed to parse output argument %s: %v", output, err)
	}

	tmp, err := ioutil.TempDir("", "tmp")
	if err != nil {
		log.Fatalf("Failed to create temporary directory: %v", err)
	}
	out, err := ioutil.TempDir("", "out")
	if err != nil {
		log.Fatalf("Failed to create temporary directory: %v", err)
	}
	return &Args{src, dst, tmp, out}
}

func withExt(path string, ext string) string {
	from := filepath.Ext(path)
	return path[0:len(path)-len(from)] + ext
}

func convert(dir string) error {
	formats := map[string]string{
		".png": ".jpg",
	}
	filepath.Walk(dir, func(src string, info os.FileInfo, err error) error {
		if to, ok := formats[filepath.Ext(src)]; ok {
			dst := withExt(src, to)
			fmt.Printf("INFO: Convert %s %s\n", src, dst)

			cmd := exec.Command("convert", src, dst)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if err := cmd.Run(); err != nil {
				return fmt.Errorf("ERROR: Failed to convert: %v", err)
			}
		}
		return nil
	})
	return nil
}

func main() {
	args := readArgs()

	if _, err := kuhnuri.DownloadFile(args.src, args.tmp); err != nil {
		log.Fatalf("Failed to download %s: %v", args.src, err)
	}

	if err := convert(args.tmp); err != nil {
		log.Fatalf("Failed to convert %s: %v", args.tmp, err)
	}

	if err := kuhnuri.UploadFile(args.tmp, args.dst); err != nil {
		log.Fatalf("Failed to upload %s: %v", args.dst, err)
	}
}
