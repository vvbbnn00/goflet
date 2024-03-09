package image

import (
	"log"
	"net/url"
	"strconv"
)

type ScaleType int
type PictureFormat string

const (
	ScaleTypeFit ScaleType = iota
	ScaleTypeFill
	ScaleTypeResize
	ScaleTypeFitWidth
	ScaleTypeFitHeight
)

const (
	PictureFormatJpeg PictureFormat = "jpeg"
	PictureFormatPng  PictureFormat = "png"
	PictureFormatGif  PictureFormat = "gif"
	//PictureFormatWebp PictureFormat = "webp"
)

// ProcessParams the parameters for processing the image
type ProcessParams struct {
	Width   int           // The width of the image
	Height  int           // The height of the image
	Scale   ScaleType     // The scale type
	Quality int           // The quality of the image
	Angle   int           // The angle of the image
	Format  PictureFormat // The format of the image
}

// Print the parameters
func (i *ProcessParams) Print() {
	log.Printf("Width: %d, Height: %d, Scale: %d, Quality: %d, Angle: %d, Format: %s",
		i.Width, i.Height, i.Scale, i.Quality, i.Angle, i.Format)
}

// Dump the parameters
func (i *ProcessParams) Dump() string {
	return "w" + strconv.Itoa(i.Width) + "h" + strconv.Itoa(i.Height) + "s" +
		strconv.Itoa(int(i.Scale)) + "q" + strconv.Itoa(i.Quality) + "a" +
		strconv.Itoa(i.Angle) + "f" + string(i.Format)
}

// GetProcessParamsFromQuery get the image process parameters from the query
func GetProcessParamsFromQuery(query url.Values) *ProcessParams {
	params := &ProcessParams{}
	if width := query.Get("w"); width != "" {
		params.Width, _ = strconv.Atoi(width)
	}
	if height := query.Get("h"); height != "" {
		params.Height, _ = strconv.Atoi(height)
	}
	if scale := query.Get("s"); scale != "" {
		switch scale {
		case "fit":
			params.Scale = ScaleTypeFit
		case "fill":
			params.Scale = ScaleTypeFill
		case "resize":
			params.Scale = ScaleTypeResize
		case "fit_width":
			params.Scale = ScaleTypeFitWidth
		case "fit_height":
			params.Scale = ScaleTypeFitHeight
		}
	}
	if quality := query.Get("q"); quality != "" {
		params.Quality, _ = strconv.Atoi(quality)
		if params.Quality < 0 || params.Quality > 100 {
			params.Quality = 90
		}
		params.Quality = params.Quality / 5 * 5
	} else {
		params.Quality = 90
	}
	if angle := query.Get("a"); angle != "" {
		params.Angle, _ = strconv.Atoi(angle)
		params.Angle = params.Angle % 360
	}
	if format := query.Get("f"); format != "" {
		switch format {
		case "jpeg":
			params.Format = PictureFormatJpeg
		case "png":
			params.Format = PictureFormatPng
		case "gif":
			params.Format = PictureFormatGif
		//case "webp":
		//	params.Format = PictureFormatWebp
		default:
			params.Format = PictureFormatJpeg
		}
	} else {
		params.Format = PictureFormatJpeg
	}

	// If the format is PNG, the quality has only 2 values: 100 or 85
	if params.Format == PictureFormatPng {
		if params.Quality < 100 {
			params.Quality = 85 // Means best compression
		}
		// Otherwise, the quality is 100
	}

	// If the format is GIF, the quality is not used
	if params.Format == PictureFormatGif {
		params.Quality = 0
	}

	return params
}