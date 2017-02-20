package scpFile

import (
	"os"
	"strings"
)

// Create is used to create a file with a specific name
func Create(fn string) (*os.File, error) {
	tfn := strings.TrimSpace(fn)
	f, err := os.Create(tfn)
	if err != nil {
		return f, err
	}

	return f, nil
}
