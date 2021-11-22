package captcha

import (
	"image"
	"image/color"
)

type ImageBuf struct {
	i image.Image
	w int
	h int
}

func (i *ImageBuf) getHeight() int {
	return i.h
}

func (i *ImageBuf) getWidth() int {
	return i.w
}

func (i *ImageBuf) getRGBA(x, y int) color.RGBA64 {
	r, g, b, a := i.i.At(x, y).RGBA()
	return color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
}

func (i *ImageBuf) setRGBA(x, y int, c color.Color) {
	switch i.i.(type) {
	case *image.RGBA:
		i.i.(*image.RGBA).Set(x, y, c)
	case *image.NRGBA:
		i.i.(*image.NRGBA).Set(x, y, c)
	}
}
