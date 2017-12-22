package text

import (
	. "comanche/util"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/golang/freetype/truetype"
)

var font *truetype.Font

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

	p := gl.CreateProgram()
	AttachShader(p, gl.VERTEX_SHADER, "text/vert.glsl")
	AttachShader(p, gl.FRAGMENT_SHADER, "text/frag.glsl")
	gl.LinkProgram(p)

	gl.Uniform1i(gl.GetUniformLocation(p, gl.Str("tex\x00")), 0)

	/*var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 4, gl.FLOAT, false, 0, nil)

	return &ShaderProgram{p, gl.GetUniformLocation(p, gl.Str("mvp\x00"))}*/
}
