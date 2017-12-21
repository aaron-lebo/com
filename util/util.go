package util

import (
	"io/ioutil"
	"path/filepath"
)

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

func Read(path string) []byte {
	abs, err := filepath.Abs("./")
	Check(err)
	file, err := ioutil.ReadFile(abs + path)
	Check(err)
	return file
}
