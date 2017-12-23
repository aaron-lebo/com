package block

import (
	. "comanche/util"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

var (
	chunk             = make([]float32, 0, 24*4096)
	chunkIndices      = make([]uint16, 0, 36*4096)
	vbo, ibo, program uint32
)

func Add(x, y, z float32) {
	const a = 0.5
	const b = -a
	block := []float32{
		b, b, a,
		a, b, a,
		a, a, a,
		b, a, a,
		b, b, b,
		b, a, b,
		a, a, b,
		a, b, b,
	}
	blockIndices := []uint16{
		0, 1, 2, 2, 3, 0, // +z
		4, 5, 6, 6, 7, 4, // -z
		3, 2, 6, 6, 5, 3, // +y
		0, 4, 7, 7, 1, 0, // -y
		1, 7, 6, 6, 2, 1, // +x
		0, 3, 5, 5, 4, 0, // -x
	}
	cnt := uint16(len(chunk) / 3)
	for row := 0; row < 8; row++ {
		idx := row * 3
		chunk = append(chunk, block[idx]+x, block[idx+1]+y, block[idx+2]+z)
	}
	for _, idx := range blockIndices {
		chunkIndices = append(chunkIndices, cnt+idx)
	}
}

func Init() {
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(chunk), gl.Ptr(chunk), gl.STATIC_DRAW)
	gl.GenBuffers(1, &ibo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ibo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 2*len(chunkIndices), gl.Ptr(chunkIndices), gl.STATIC_DRAW)

	program = CreateProgram("")
}

func Render(mvp mgl32.Mat4) {
	gl.UseProgram(program)
	gl.UniformMatrix4fv(gl.GetUniformLocation(program, gl.Str("mvp\x00")), 1, false, &mvp[0])
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	defer gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ibo)
	defer gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
	gl.EnableVertexAttribArray(0)
	defer gl.DisableVertexAttribArray(0)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)
	gl.DrawElements(gl.TRIANGLES, int32(len(chunkIndices)), gl.UNSIGNED_SHORT, nil)

}
