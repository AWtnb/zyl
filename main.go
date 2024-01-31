package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/AWtnb/moko/launchentry"
	"github.com/ktr0731/go-fuzzyfinder"
)

func main() {
	var (
		src     string
		filer   string
		all     bool
		exclude string
	)
	flag.StringVar(&src, "src", "", "source yaml file path")
	flag.StringVar(&filer, "filer", "explorer.exe", "filer path")
	flag.BoolVar(&all, "all", false, "switch to search including file")
	flag.StringVar(&exclude, "exclude", "", "search exception (comma-separated)")
	flag.Parse()
	os.Exit(run(src, filer, all, exclude))
}

type Filer struct {
	path string
}

func (fl *Filer) setPath(path string) {
	if _, err := os.Stat(path); err == nil {
		fl.path = path
		return
	}
	fl.path = "explorer.exe"
}

func (fl Filer) open(path string) {
	if fi, err := os.Stat(path); err == nil && fi.IsDir() {
		exec.Command(fl.path, path).Start()
		return
	}
	exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", path).Start()
}

func run(src string, filer string, all bool, exclude string) int {
	if len(src) < 1 {
		p, _ := os.Executable()
		src = filepath.Join(filepath.Dir(p), "launch.yaml")
	}

	var les launchentry.LaunchEntries
	if err := les.Load(src); err != nil {
		fmt.Println(err)
		return 1
	}

	selected, err := les.Select()
	if err != nil {
		if err != fuzzyfinder.ErrAbort {
			fmt.Println(err)
		}
		return 1
	}

	var t launchentry.Target
	t.SetEntry(selected)
	if !t.IsValid() {
		fmt.Printf("invalid path: '%s'\n", selected.Path)
		return 1
	}

	var fl Filer
	fl.setPath(filer)

	if t.IsUri() || t.IsFile() {
		fl.open(t.Path())
		return 0
	}

	cs, err := t.GetChildItem(all, exclude)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	if len(cs) < 1 {
		fl.open(t.Path())
		return 0
	}
	c, err := t.SelectItem(cs)
	if err != nil {
		if err != fuzzyfinder.ErrAbort {
			fmt.Println(err)
		}
		return 1
	}
	fl.open(c)
	return 0
}
