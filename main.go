package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/AWtnb/zyl/launchentry"
	"github.com/ktr0731/go-fuzzyfinder"
)

func main() {
	var (
		src     string
		filer   string
		all     bool
		exclude string
		stdout  bool
	)
	flag.StringVar(&src, "src", "", "source yaml file path")
	flag.StringVar(&filer, "filer", "explorer.exe", "filer path")
	flag.BoolVar(&all, "all", false, "switch to search including file")
	flag.StringVar(&exclude, "exclude", "", "search exception (comma-separated)")
	flag.BoolVar(&stdout, "stdout", false, "switch to stdout")
	flag.Parse()

	var f Filer
	f.Init(filer)

	if len(src) < 1 {
		p, _ := os.Executable()
		src = filepath.Join(filepath.Dir(p), "launch.yaml")
	}
	os.Exit(run(src, f, all, exclude, stdout))
}

func find(src string, all bool, exclude string) (string, error) {
	var les launchentry.LaunchEntries
	if err := les.Load(src); err != nil {
		return "", err
	}

	sel, err := les.Select()
	if err != nil {
		return "", err
	}

	var t launchentry.Target
	t.SetEntry(sel)
	if t.IsInvalid() {
		return sel.Path, fmt.Errorf("invalid path: %s", sel.Path)
	}

	if t.IsFile() || t.IsUri() {
		return t.Path(), nil
	}

	evr, fd, err := t.GetChildItem(all, exclude)
	if err != nil {
		return "", err
	}
	if len(fd) < 1 {
		return t.Path(), nil
	}
	var prompt string
	if evr {
		prompt = "#"
	} else {
		prompt = ">"
	}
	c, err := t.SelectItem(fd, prompt)
	if err != nil {
		return "", err
	}
	return c, nil
}

func run(src string, flr Filer, all bool, exclude string, stdout bool) int {
	p, err := find(src, all, exclude)
	if err != nil {
		if err != fuzzyfinder.ErrAbort {
			fmt.Println(err.Error())
		}
		return 1
	}
	if stdout {
		fmt.Println(p)
		return 0
	}

	if err := flr.OpenSmart(p, ""); err != nil {
		fmt.Println(err.Error())
		return 1
	}
	return 0
}
