package nds

import (
	"bytes"
	"encoding/binary"
	"github.com/szerookii/desmume-launcher/utils/imageUtils"
	"image"
	"image/color"
	"strings"

	"github.com/fogleman/gg"
)

type NDSBanner struct {
	Version       uint16
	Crc16         uint16
	Reserved      [28]byte
	TileData      [512]byte
	Palette       [32]byte
	JapaneseTitle [256]byte
	EnglishTitle  [256]byte
	FrenchTitle   [256]byte
	GermanTitle   [256]byte
	ItalianTitle  [256]byte
	SpanishTitle  [256]byte
}

func (banner *NDSBanner) IconPNG(imageSize int) ([]byte, error) {
	palette := make([]color.Color, 16)

	for i := 0; i < 16; i++ {
		colorIndex := binary.LittleEndian.Uint16(banner.Palette[i*2 : (i+1)*2])
		palette[i] = color.RGBA{
			R: uint8((colorIndex & 0x001F) << 3),
			G: uint8((colorIndex & 0x03E0) >> 2),
			B: uint8((colorIndex & 0x7C00) >> 7),
			A: 0xFF,
		}
	}

	icon := image.NewRGBA(image.Rect(0, 0, 32, 32))

	for tI := 0; tI < 16; tI++ {
		for pI := 0; pI < 32; pI++ {
			pixelData := banner.TileData[(tI<<5)+pI]

			tileX := ((tI & 3) << 3) + ((pI << 1) & 7)
			tileY := ((tI >> 2) << 3) + (pI >> 2)

			upNibble, lowNibble := pixelData>>4, pixelData&0x0F

			if upNibble != 0 {
				icon.Set(tileX+1, tileY, palette[upNibble])
			}
			if lowNibble != 0 {
				icon.Set(tileX, tileY, palette[lowNibble])
			}
		}
	}

	dc := gg.NewContext(imageSize, imageSize)
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	dc.DrawImageAnchored(imageUtils.ResizeImage(icon, imageSize), imageSize/2, imageSize/2, 0.5, 0.5)

	imageBuffer := bytes.NewBuffer(nil)
	if err := dc.EncodePNG(imageBuffer); err != nil {
		return nil, err
	}

	return imageBuffer.Bytes(), nil
}

func (banner *NDSBanner) Author() string {
	cleanedText := CleanText(ISO2UTF8(string(banner.EnglishTitle[:])))

	return strings.Split(cleanedText, " ")[len(strings.Split(cleanedText, " "))-1]
}

func (banner *NDSBanner) JapaneseTitleString() string {
	cleanedText := CleanText(ISO2UTF8(string(banner.JapaneseTitle[:])))

	return strings.Join(strings.Split(cleanedText, " ")[:len(strings.Split(cleanedText, " "))-1], " ")
}

func (banner *NDSBanner) EnglishTitleString() string {
	cleanedText := ISO2UTF8(CleanText(string(banner.EnglishTitle[:])))

	return strings.Join(strings.Split(cleanedText, " ")[:len(strings.Split(cleanedText, " "))-1], " ")
}

func (banner *NDSBanner) FrenchTitleString() string {
	cleanedText := CleanText(ISO2UTF8(string(banner.FrenchTitle[:])))

	return strings.Join(strings.Split(cleanedText, " ")[:len(strings.Split(cleanedText, " "))-1], " ")
}

func (banner *NDSBanner) GermanTitleString() string {
	cleanedText := CleanText(ISO2UTF8(string(banner.GermanTitle[:])))

	return strings.Join(strings.Split(cleanedText, " ")[:len(strings.Split(cleanedText, " "))-1], " ")
}

func (banner *NDSBanner) ItalianTitleString() string {
	cleanedText := CleanText(ISO2UTF8(string(banner.ItalianTitle[:])))

	return strings.Join(strings.Split(cleanedText, " ")[:len(strings.Split(cleanedText, " "))-1], " ")
}

func (banner *NDSBanner) SpanishTitleString() string {
	cleanedText := CleanText(ISO2UTF8(string(banner.SpanishTitle[:])))

	return strings.Join(strings.Split(cleanedText, " ")[:len(strings.Split(cleanedText, " "))-1], " ")
}
