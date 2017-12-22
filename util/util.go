package util

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

func ReadFile(path string) []byte {
	abs, err := filepath.Abs("./")
	Check(err)
	file, err := ioutil.ReadFile(abs + "/" + path)
	Check(err)
	return file
}

func AttachShader(program, kind uint32, path string) {
	s := gl.CreateShader(kind)
	strs, free := gl.Strs(string(ReadFile(path)) + "\x00")
	gl.ShaderSource(s, 1, strs, nil)
	free()
	gl.CompileShader(s)
	gl.AttachShader(program, s)

	var status int32
	gl.GetShaderiv(s, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var length int32
		gl.GetShaderiv(s, gl.INFO_LOG_LENGTH, &length)
		log := strings.Repeat("\x00", int(length+1))
		gl.GetShaderInfoLog(s, length, nil, gl.Str(log))
		panic(fmt.Sprintf("shader compilation failed: %v\n%v", log, path))
	}
}
