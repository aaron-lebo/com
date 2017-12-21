package text

import (
	. "comanche/util"
	"github.com/golang/freetype/truetype"
)

var font *truetype.Font

func init() {
	ttf := ReadFile("text/NotoMono-Regular.ttf")
	font, err := truetype.Parse(ttf)
	Check(err)
	println(font)
}

func Print() {
	println(font)
}
