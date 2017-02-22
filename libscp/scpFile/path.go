package scpFile

import (
	"os"
	"path"
)

// ExpandPath is used to ensure a path that we have is fully expanded rather
// than something like ./LICENSE
func ExpandPath(f string) (string, error) {
	if !path.IsAbs(f) {
		wd, err := os.Getwd()
		if err != nil {
			return "", err
		}
		fname := path.Clean(path.Join(wd, f))
		return fname, nil
	}

	return f, nil
}
