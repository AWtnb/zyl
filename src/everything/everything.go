// https://github.com/jof4002/Everything/blob/master/_Example/walk/example.go
package everything

import (
	"os"
	"path/filepath"

	"github.com/AWtnb/moko/everything/core"
)

func Scan(query string, skipFile bool) []string {
	sl := []string{}
	if err := checkDll("Everything64.dll"); err != nil {
		return sl
	}
	core.Walk(query, skipFile, func(path string, isFile bool) error {
		if skipFile && isFile {
			return nil
		}
		sl = append(sl, path)
		return nil
	})
	return sl
}

func getExeDir() string {
	if exePath, err := os.Executable(); err != nil {
		return exePath
	}
	return ""
}

func checkDll(name string) error {
	exeDir := getExeDir()
	path := filepath.Join(exeDir, name)
	_, err := os.Stat(path)
	return err
}
