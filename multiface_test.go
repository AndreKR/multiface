package multiface_test

import (
	"testing"
	"golang.org/x/image/font"
	"image"
	"github.com/golang/freetype/truetype"
	"github.com/golang/freetype"
	"io/ioutil"
	"github.com/AndreKR/multiface"
	"image/draw"
	"os"
	"image/png"
	"bytes"
)

func TestMultiface(t *testing.T) {

	face := new(multiface.Face)

	opts := &truetype.Options{Size: 20, DPI: 96}

	var fnt *truetype.Font
	var fc font.Face

	// Add ArchitectsDaughter font, which does not include a glyph for ก, but has a handwriting-style glyph for a
	fnt = readFont(t, "testdata/ArchitectsDaughter-Regular.ttf")
	fc = truetype.NewFace(fnt, opts)
	face.AddTruetypeFace(fc, fnt)

	// Add Kanit font, which does have a glyph for ก
	fnt = readFont(t, "testdata/Kanit-Regular.ttf")
	fc = truetype.NewFace(fnt, opts)
	face.AddTruetypeFace(fc, fnt)

	img := image.NewRGBA(image.Rect(0, 0, 50, 50))
	draw.Draw(img, img.Rect, image.White, image.ZP, draw.Src)

	d := font.Drawer{}
	d.Dst = img
	d.Src = image.Black
	d.Face = face
	d.Dot = freetype.Pt(10, 25)
	d.DrawString("กa")

	f, err := os.Create("testdata/output.png")
	checkErr(t, err)
	err = png.Encode(f, img)
	checkErr(t, err)
	err = f.Close()
	checkErr(t, err)

	f, err = os.Open("testdata/reference.png")
	checkErr(t, err)
	ref, err  := png.Decode(f)
	checkErr(t, err)

	if bytes.Compare(ref.(*image.RGBA).Pix, img.Pix) != 0 {
		t.Fatal("output does not match reference")
	}
}

func readFont(t *testing.T, filename string) *truetype.Font {
	data, err := ioutil.ReadFile(filename)
	checkErr(t, err)
	fnt, err := truetype.Parse(data)
	checkErr(t, err)
	return fnt
}

func checkErr(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}