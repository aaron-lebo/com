package text

import (
	. "comanche/util"
	"github.com/go-gl/gl/v4.1-core/gl"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/font/inconsolata"
	"golang.org/x/image/math/fixed"
	"image"
	"image/draw"
)

var (
	face         *basicfont.Face
	vbo, program uint32
	attr_pos     uint32
)

func init() {
	face = inconsolata.Regular8x16
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

	rgba := image.NewRGBA(face.Mask.Bounds())
	draw.Draw(rgba, rgba.Bounds(), face.Mask, image.Point{0, 0}, draw.Src)
	size := rgba.Rect.Size()
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(size.X), int32(size.Y), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba.Pix))

	gl.GenBuffers(1, &vbo)

	program = CreateProgram("test/")
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

	const sx = float32(2.0 / 800)
	const sy = float32(2.0 / 600)
	w := float32(face.Width) * sx
	h := float32(face.Height) * sy
	n := 6 * len(text)
	coords := make([]float32, 0, n)
	for _, chr := range text {
		x2 := x + float32(face.Left)*sx
		y2 := -y - float32(face.Ascent)*sy
		x += w
		_, _, pos, _, _ := face.Glyph(fixed.Point26_6{0, 0}, chr)
		txy := float32(pos.Y)
		coords = append(coords,
			x2, -y2-h, 0, txy+h,
			x2+w, -y2-h, w, txy+h,
			x2+w, -y2, w, txy,
			x2+w, -y2, w, txy,
			x2, -y2, 0, txy,
			x2, -y2-h, 0, txy+h,
		)

	}
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(coords), gl.Ptr(coords), gl.DYNAMIC_DRAW)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(n))
}
