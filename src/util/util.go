package util

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func HasFile(path string) bool {
	nf := 0
	items, err := ioutil.ReadDir(path)
	if err == nil {
		for _, item := range items {
			if !item.IsDir() {
				nf++
			}
		}
	}
	return nf > 0
}

func ExecuteFile(path string) {
	exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", path).Start()
}

func IsValidPath(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func IsExecutable(path string) bool {
	if strings.HasPrefix(path, "http") {
		return true
	}
	fi, _ := os.Stat(path)
	return !fi.IsDir()
}

func ToSlice(s string, sep string) []string {
	var ss []string
	for _, elem := range strings.Split(s, sep) {
		ss = append(ss, strings.TrimSpace(elem))
	}
	return ss
}
