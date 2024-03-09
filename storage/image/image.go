package image

import (
	"bytes"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/nfnt/resize"
	"github.com/vvbbnn00/goflet/config"
	"github.com/vvbbnn00/goflet/util/log"
	"image"
	"image/color"
	"os"
)

// ProcessImage process the image with the given parameters
func ProcessImage(fs *os.File, p *ProcessParams) (*bytes.Buffer, error) {
	conf := config.GofletCfg.ImageConfig
	decoded, _, err := image.Decode(fs)
	if err != nil {
		return nil, err
	}

	width, height := decoded.Bounds().Dx(), decoded.Bounds().Dy()
	if width > conf.MaxWidth || height > conf.MaxHeight {
		return nil, fmt.Errorf("image size is too large")
	}

	// Resize the image
	resized := resizeImage(decoded, p.Scale, p.Width, p.Height)

	// Rotate the image
	rotated := rotateImage(resized, p.Angle)

	// Change the format
	buf, err := convertImageFormat(rotated, p.Format, p.Quality)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

// resizeImage resize the image with the given parameters
func resizeImage(img image.Image, scaleType ScaleType, width, height int) image.Image {
	imgWidth, imgHeight := img.Bounds().Dx(), img.Bounds().Dy()

	// If the width and height are 0, return the original image
	if width == 0 && height == 0 {
		return img
	}
	// Ensure the image is not smaller than the requested size
	if width > imgWidth || height > imgHeight {
		return img
	}

	// If the width or height is 0, calculate the aspect ratio
	if width == 0 {
		width = imgWidth * height / imgHeight
	}
	if height == 0 {
		height = imgHeight * width / imgWidth
	}

	switch scaleType {
	case ScaleTypeFit:
		widthRatio := float64(width) / float64(imgWidth)
		heightRatio := float64(height) / float64(imgHeight)
		if widthRatio < heightRatio {
			height = int(float64(imgHeight) * widthRatio)
		} else {
			width = int(float64(imgWidth) * heightRatio)
		}
	case ScaleTypeFill:
		widthRatio := float64(width) / float64(imgWidth)
		heightRatio := float64(height) / float64(imgHeight)
		if widthRatio > heightRatio {
			height = int(float64(imgHeight) * widthRatio)
		} else {
			width = int(float64(imgWidth) * heightRatio)
		}
	case ScaleTypeFitWidth:
		height = int(float64(imgHeight) * float64(width) / float64(imgWidth))
	case ScaleTypeFitHeight:
		width = int(float64(imgWidth) * float64(height) / float64(imgHeight))
	case ScaleTypeResize:
		// Do nothing
	}

	return resize.Resize(uint(width), uint(height), img, resize.Lanczos3)
}

// rotateImage rotate the image with the given angle
func rotateImage(img image.Image, angle int) image.Image {
	if angle%360 == 0 {
		return img
	}
	angle = angle % 360 // Normalize the angle
	log.Warnf("Rotating the image by %d degrees", angle)
	return imaging.Rotate(img, float64(angle), color.Transparent)
}
