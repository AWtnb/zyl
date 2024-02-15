// https://github.com/jof4002/Everything/blob/master/_Example/walk/example.go
package everything

func Scan(query string, skipFile bool) (sl []string, err error) {
	err = Walk(query, skipFile, func(path string, isFile bool) error {
		if skipFile && isFile {
			return nil
		}
		sl = append(sl, path)
		return nil
	})
	return
}
