package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	pdf "github.com/unidoc/unipdf/v3/model"
)

var cfg struct {
	files []string
	out   string
}

func init() {
	if len(os.Args) < 3 {
		fatalf("Usage: %s OUTPUT FILE [FILE...]\n", os.Args[0])
	}
	cfg.out = os.Args[1]
	for _, pattern := range os.Args[2:] {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			fatalf("filepath.Glob: %v\n", err)
		}
		for _, match := range matches {
			cfg.files = append(cfg.files, match)
		}
	}
	if len(cfg.files) <= 1 {
		fatalf("Not enough files to merge.\n")
	}
}

func main() {
	out, err := os.Create(cfg.out)
	if err != nil {
		fatalf("os.Create: %v\n", err)
	}
	defer out.Close()
	if err := merge(out, cfg.files); err != nil {
		fatalf("merge: %v", err)
	}
}

func merge(out io.Writer, files []string) error {
	wr := pdf.NewPdfWriter()

	for _, file := range files {
		if err := include(&wr, file); err != nil {
			return err
		}
	}

	return wr.Write(out)
}

func include(wr *pdf.PdfWriter, file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	r, err := pdf.NewPdfReader(f)
	if err != nil {
		return err
	}

	pages, err := r.GetNumPages()
	if err != nil {
		return err
	}

	for i := 0; i < pages; i++ {
		page, err := r.GetPage(i + 1)
		if err != nil {
			return err
		}

		if err := wr.AddPage(page); err != nil {
			return err
		}
	}

	return nil
}

func errorf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
}

func fatalf(format string, args ...interface{}) {
	errorf(format, args...)
	os.Exit(1)
}
