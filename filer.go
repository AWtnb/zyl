package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
	fmt.Printf("open '%s' itself? (y/N): ", filepath.Base(path))
	sc := bufio.NewScanner(os.Stdin)
	sc.Scan()
	s := sc.Text()
	if strings.EqualFold(s, "y") {
		fmt.Println("[ACCEPTED] opening in default app.")
		exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", path).Start()
		return
	}
	fmt.Println("[DENIED] opening its directory on filer.")
	d := filepath.Dir(path)
	flr.Open(d)
}
