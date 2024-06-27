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
	)
	flag.StringVar(&src, "src", "", "source yaml file path")
	flag.StringVar(&filer, "filer", "explorer.exe", "filer path")
	flag.BoolVar(&all, "all", false, "switch to search including file")
	flag.StringVar(&exclude, "exclude", "", "search exception (comma-separated)")
	flag.Parse()

	var f Filer
	f.SetPath(filer)
	if len(src) < 1 {
		p, _ := os.Executable()
		src = filepath.Join(filepath.Dir(p), "launch.yaml")
	}
	os.Exit(run(src, f, all, exclude))
}

func run(src string, flr Filer, all bool, exclude string) int {
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
	if t.IsInvalid() {
		fmt.Printf("invalid path: '%s'\n", selected.Path)
		return 1
	}

	if err := t.RunApp(); err == nil {
		return 0
	}

	withEv, cs, err := t.GetChildItem(all, exclude)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	if len(cs) < 1 {
		flr.Open(t.Path())
		return 0
	}
	var prompt string
	if withEv {
		prompt = "#"
	} else {
		prompt = ">"
	}
	c, err := t.SelectItem(cs, prompt)
	if err != nil {
		if err != fuzzyfinder.ErrAbort {
			fmt.Println(err)
		}
		return 1
	}
	flr.OpenSmart(c, "")
	return 0
}
