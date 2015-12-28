package goScp

import (
	"os"
	"strings"
	"log"
)

func createNewFile(filename string) *os.File {
	file, err := os.Create(strings.TrimSpace(filename))
	if err != nil {
		log.Fatal(err)
	}

	return file
}

func writeParitalToFile(file *os.File, content []byte) {
	_, err := file.Write(content)
	if err != nil {
		log.Fatal(err)
	}
}
