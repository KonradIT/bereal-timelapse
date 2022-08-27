package cmd

import (
	"image"
	"image/draw"
	"image/jpeg"
	"os"

	"github.com/nfnt/resize"
)

const (
	PADDING_TOP_LEFT   = 30
	HEADSHOT_RESIZE_TO = 600
)

func Superimpose(
	background string,
	foreground string,
	output string,
) error {

	bgImage, err := os.Open(background)
	if err != nil {
		return err
	}

	bg, err := jpeg.Decode(bgImage)
	if err != nil {
		return err
	}
	defer bgImage.Close()

	fgImage, err := os.Open(foreground)
	if err != nil {
		return err
	}
	fg, err := jpeg.Decode(fgImage)
	if err != nil {
		return err
	}
	defer fgImage.Close()

	offset := image.Pt(PADDING_TOP_LEFT, PADDING_TOP_LEFT)
	b := bg.Bounds()
	mixedImage := image.NewRGBA(b)

	newImage := resize.Resize(HEADSHOT_RESIZE_TO, 0, fg, resize.Lanczos3)

	draw.Draw(mixedImage, b, bg, image.ZP, draw.Src)
	draw.Draw(mixedImage, newImage.Bounds().Add(offset), newImage, image.ZP, draw.Over)

	render, err := os.Create(output)
	if err != nil {
		return err
	}
	err = jpeg.Encode(render, mixedImage, &jpeg.Options{
		Quality: jpeg.DefaultQuality,
	})
	if err != nil {
		return err
	}
	defer render.Close()

	return nil
}
