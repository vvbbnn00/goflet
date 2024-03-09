package image

import (
	"bytes"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"

	"github.com/pkg/errors"
)

// convertImageFormat convert the image to the given format
func convertImageFormat(img image.Image, format PictureFormat, quality int) (*bytes.Buffer, error) {
	var buf bytes.Buffer

	switch format {
	case PictureFormatJpeg:
		if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality}); err != nil {
			return nil, err
		}
	case PictureFormatPng:
		compression := png.DefaultCompression
		if quality != 100 {
			compression = png.BestCompression
		}
		encoder := png.Encoder{
			CompressionLevel: compression,
		}
		if err := encoder.Encode(&buf, img); err != nil {
			return nil, err
		}
	case PictureFormatGif:
		if err := gif.Encode(&buf, img, nil); err != nil {
			return nil, err
		}
	// case PictureFormatWebp:
	//	if err := webp.Encode(&buf, img, &webp.Options{Lossless: false, Quality: float32(quality)}); err != nil {
	//		return nil, err
	//	}
	default:
		return nil, errors.Errorf("unsupported format: %s", format)
	}

	return &buf, nil
}
