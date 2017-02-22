package scpFile

import "os"

// WriteBytes is used to write an array of bytes to a file
func WriteBytes(file *os.File, content []byte) (int, error) {
	w, err := file.Write(content)
	if err != nil {
		return 0, err
	}

	return w, nil
}
