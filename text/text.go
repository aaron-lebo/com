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
	sizeX, sizeY int
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
	sizeX = size.X
	sizeY = size.Y
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(sizeX), int32(sizeY), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba.Pix))

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

	per := float32(17) / float32(sizeY)
	const sx = float32(2.0 / 800)
	const sy = float32(2.0 / 600)
	w := float32(face.Width) * sx
	h := float32(17) * sy
	n := 6 * len(text)
	coords := make([]float32, 0, n)
	for _, chr := range text {
		x2 := x + float32(face.Left)*sx
		y2 := -y - float32(face.Ascent)*sy
		x += w
		_, _, pos, _, _ := face.Glyph(fixed.Point26_6{100, 100}, chr)
		i := float32(pos.Y / 17)
		i1 := i + 1.0
		coords = append(coords,
			x2, -y2-h, 0, per*i1,
			x2+w, -y2-h, 1, per*i1,
			x2+w, -y2, 1, per*i,
			x2+w, -y2, 1, per*i,
			x2, -y2, 0, per*i,
			x2, -y2-h, 0, per*i1,
		)
	}
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(coords), gl.Ptr(coords), gl.DYNAMIC_DRAW)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(n))
}
