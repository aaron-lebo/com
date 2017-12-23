package text

import (
	. "comanche/util"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/golang/freetype/truetype"
)

var (
	font         *truetype.Font
	vbo, program uint32
	attr_pos     uint32
)

func init() {
	ttf := ReadFile("text/NotoMono-Regular.ttf")
	var err error
	font, err = truetype.Parse(ttf)
	Check(err)
}

func Init() {
	gl.ActiveTexture(gl.TEXTURE0)
	var tex uint32
	gl.GenTextures(1, &tex)
	gl.BindTexture(gl.TEXTURE_2D, tex)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	gl.GenBuffers(1, &vbo)

	program = CreateProgram("text/")
	attr_pos = uint32(gl.GetAttribLocation(program, gl.Str("pos\x00")))
	gl.Uniform1i(gl.GetUniformLocation(program, gl.Str("tex\x00")), 0)
}

func Render(text string, x, y float32) {
	gl.UseProgram(program)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	defer gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.EnableVertexAttribArray(attr_pos)
	defer gl.DisableVertexAttribArray(attr_pos)

	gl.VertexAttribPointer(attr_pos, 4, gl.FLOAT, false, 0, nil)
	for _, chr := range text {
		var g truetype.GlyphBuf
		Check(g.Load(font, 12, font.Index(chr), 0))

		max := g.Bounds.Max
		w := float32(max.X)
		h := float32(max.Y)
		gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RED, int32(w), int32(h), 0, gl.RED, gl.UNSIGNED_BYTE, gl.Ptr(g.Points))
		min := g.Bounds.Min
		x2 := x + float32(min.X)*12
		y2 := -y + float32(min.Y)*12
		w *= 12
		h *= 12
		box := []float32{
			x2, -y2, 0, 0,
			x2 + w, -y2, 1, 0,
			x2, -y2 - h, 0, 1,
			x2 + w, -y2 - h, 1, 1,
		}
		gl.BufferData(gl.ARRAY_BUFFER, 4*len(box), gl.Ptr(box), gl.DYNAMIC_DRAW)
		gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
	}
}
