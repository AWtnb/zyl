package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/AWtnb/moko/sys"
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

func (flr Filer) Open(path string) {
	exec.Command(flr.path, path).Start()
}

func (flr Filer) OpenSmart(path string) {
	if fi, err := os.Stat(path); err == nil && fi.IsDir() {
		flr.Open(path)
		return
	}
	fmt.Printf("'%s' is a file.\nopen itself? (y/N): ", path)
	sc := bufio.NewScanner(os.Stdin)
	sc.Scan()
	s := sc.Text()
	if strings.EqualFold(s, "y") {
		sys.Open(path)
		fmt.Printf("[Y] opening in default app: '%s'\n", path)
		return
	}
	d := filepath.Dir(path)
	flr.Open(d)
	fmt.Printf("[N] opening its directory: '%s'\n", d)
}
