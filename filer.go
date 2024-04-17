package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/AWtnb/moko/sys"
	"github.com/ktr0731/go-fuzzyfinder"
)

type Filer struct {
	path string
}

func (flr *Filer) SetPath(path string) {
	if _, err := os.Stat(path); err == nil {
		flr.path = path
		return
	}
	flr.path = "explorer.exe"
}

func (flr Filer) Open(path string) error {
	if fi, err := os.Stat(path); err == nil && fi.IsDir() {
		exec.Command(flr.path, path).Start()
		return nil
	}
	return fmt.Errorf("failed to open '%s' on filer", path)
}

func (flr Filer) OpenSmart(path string, curDir string) {
	if err := flr.Open(path); err == nil {
		return
	}
	if filepath.Dir(path) == curDir {
		sys.Open(path)
		return
	}
	d := filepath.Dir(path)
	ss := []string{path, d}
	idx, err := fuzzyfinder.Find(ss, func(i int) string {
		p := ss[i]
		rel, _ := filepath.Rel(filepath.Dir(d), p)
		return rel
	}, fuzzyfinder.WithCursorPosition(fuzzyfinder.CursorPositionTop))
	if err != nil {
		return
	}

	p := ss[idx]
	if err := flr.Open(p); err != nil {
		sys.Open(p)
	}
}
