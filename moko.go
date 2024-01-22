package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/AWtnb/moko/launchentry"
	"github.com/AWtnb/moko/selectedentry"
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

	var se selectedentry.SelectedEntry
	se.SetEntry(selected)
	if !se.IsValid() {
		fmt.Printf("invalid path: '%s'\n", selected.Path)
		return 1
	}

	se.SetFiler(filer)

	if se.IsExecutable() {
		se.OpenSelf()
		return 0
	}

	cs, err := se.GetChildItem(all, exclude)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	if len(cs) < 1 {
		se.OpenSelf()
		return 0
	}
	c, err := se.SelectItem(cs)
	if err != nil {
		if err != fuzzyfinder.ErrAbort {
			fmt.Println(err)
		}
		return 1
	}
	se.Open(c)
	return 0
}
