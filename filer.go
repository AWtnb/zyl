package main

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ktr0731/go-fuzzyfinder"
)

func defaultOpen(path string) error {
	return exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", path).Start()
}

type Filer struct {
	path string
}

func (flr *Filer) SetPath(path string) {
	if fi, err := os.Stat(path); err == nil && !fi.IsDir() {
		flr.path = path
		return
	}
	flr.path = "explorer.exe"
}

func (flr Filer) open(path string) error {
	return exec.Command(flr.path, path).Start()
}

func (flr Filer) OpenSmart(path string, curDir string) error {
	if fi, err := os.Stat(path); err == nil && fi.IsDir() {
		return flr.open(path)
	}
	if filepath.Dir(path) == curDir {
		return defaultOpen(path)
	}
	d := filepath.Dir(path)
	ss := []string{path, d}
	idx, err := fuzzyfinder.Find(ss, func(i int) string {
		p := ss[i]
		rel, _ := filepath.Rel(filepath.Dir(d), p)
		return filepath.ToSlash(rel)
	}, fuzzyfinder.WithCursorPosition(fuzzyfinder.CursorPositionTop))
	if err != nil {
		return err
	}
	p := ss[idx]
	fi, err := os.Stat(p)
	if err != nil {
		return err
	}
	if fi.IsDir() {
		return flr.open(path)
	}
	return defaultOpen(p)
}
