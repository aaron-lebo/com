package text

import (
	. "comanche/util"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/golang/freetype/truetype"
)

var (
	font         *truetype.Font
	vbo, program uint32
)

func init() {
	ttf := ReadFile("text/NotoMono-Regular.ttf")
	var err error
	font, err = truetype.Parse(ttf)
	Check(err)
}

func Init() {
	var tex uint32
	gl.ActiveTexture(gl.TEXTURE0)
	gl.GenTextures(1, &tex)
	gl.BindTexture(gl.TEXTURE_2D, tex)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	gl.GenBuffers(1, &vbo)

	program = CreateProgram("text/")
	gl.Uniform1i(gl.GetUniformLocation(program, gl.Str("tex\x00")), 0)
}

func Render(text string, x, y float32) {
	gl.UseProgram(program)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	pos := uint32(gl.GetAttribLocation(program, gl.Str("pos\x00")))
	gl.EnableVertexAttribArray(pos)
	//gl.VertexAttribPointer(pos, 4, gl.FLOAT, false, 0, nil)
	gl.DisableVertexAttribArray(pos)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}
