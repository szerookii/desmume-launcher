package imageUtils

import "image"

func ResizeImage(img image.Image, newSize int) *image.RGBA {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	newImg := image.NewRGBA(image.Rect(0, 0, newSize, newSize))

	for x := 0; x < newSize; x++ {
		for y := 0; y < newSize; y++ {
			srcX := x * width / newSize
			srcY := y * height / newSize

			pixelColor := img.At(srcX, srcY)

			newImg.Set(x, y, pixelColor)
		}
	}

	return newImg
}
