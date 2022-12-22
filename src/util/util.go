package util

import (
	"os"
	"os/exec"
	"strings"
)

func ExecuteFile(path string) {
	exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", path).Start()
}

func IsValidPath(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func IsDir(path string) bool {
	if fi, err := os.Stat(path); err == nil && fi.IsDir() {
		return true
	}
	return false
}

func ToSlice(s string, sep string) []string {
	var ss []string
	for _, elem := range strings.Split(s, sep) {
		ss = append(ss, strings.TrimSpace(elem))
	}
	return ss
}

var envMap = map[string]string{
	"%APPDATA%":     os.Getenv("APPDATA"),
	"%USERNAME%":    os.Getenv("USERNAME"),
	"%USERPROFILE%": os.Getenv("USERPROFILE"),
}

func ResolveEnvPath(s string) string {
	for k, v := range envMap {
		s = strings.ReplaceAll(s, k, v)
	}
	return s
}
