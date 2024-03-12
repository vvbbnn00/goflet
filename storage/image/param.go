package image

import (
	"net/url"
	"strconv"

	"github.com/vvbbnn00/goflet/config"
	"github.com/vvbbnn00/goflet/util/log"
)

// ScaleType the type of the scale
type ScaleType int

// PictureFormat the format of the picture
type PictureFormat string

const (
	// ScaleTypeFit The fit scale type
	ScaleTypeFit ScaleType = iota
	// ScaleTypeFill The fill scale type
	ScaleTypeFill
	// ScaleTypeResize The resize scale type
	ScaleTypeResize
	// ScaleTypeFitWidth The fit width scale type
	ScaleTypeFitWidth
	// ScaleTypeFitHeight The fit height scale type
	ScaleTypeFitHeight
)

const (
	// PictureFormatJpeg The jpeg picture format
	PictureFormatJpeg PictureFormat = "jpeg"
	// PictureFormatPng The png picture format
	PictureFormatPng PictureFormat = "png"
	// PictureFormatGif The gif picture format
	PictureFormatGif PictureFormat = "gif"
	// PictureFormatWebp PictureFormat = "webp"
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
	log.Debugf("Width: %d, Height: %d, Scale: %d, Quality: %d, Angle: %d, Format: %s",
		i.Width, i.Height, i.Scale, i.Quality, i.Angle, i.Format)
}

// Dump the parameters
func (i *ProcessParams) Dump() string {
	return "w" + strconv.Itoa(i.Width) + "h" + strconv.Itoa(i.Height) + "s" +
		strconv.Itoa(int(i.Scale)) + "q" + strconv.Itoa(i.Quality) + "a" +
		strconv.Itoa(i.Angle) + "f" + string(i.Format)
}

// whether the int in the array
func in(num int, arr []int) bool {
	for _, v := range arr {
		if v == num {
			return true
		}
	}
	return false
}

// GetProcessParamsFromQuery get the image process parameters from the query
func GetProcessParamsFromQuery(query url.Values) *ProcessParams {
	conf := config.GofletCfg.ImageConfig

	params := &ProcessParams{}
	if width := query.Get("w"); width != "" {
		params.Width, _ = strconv.Atoi(width)
	}
	if *conf.StrictMode && !in(params.Width, conf.AllowedSizes) {
		params.Width = 0
	} else {
		params.Width = max(params.Width, 0)
	}

	if height := query.Get("h"); height != "" {
		params.Height, _ = strconv.Atoi(height)
	}
	if *conf.StrictMode && !in(params.Height, conf.AllowedSizes) {
		params.Height = 0
	} else {
		params.Height = max(params.Height, 0)
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
		params.Angle %= 360
	}
	if format := query.Get("f"); format != "" {
		switch format {
		case "jpeg":
			params.Format = PictureFormatJpeg
		case "png":
			params.Format = PictureFormatPng
		case "gif":
			params.Format = PictureFormatGif
		// case "webp":
		//	params.Format = PictureFormatWebp
		default:
			params.Format = PictureFormat(conf.DefaultFormat)
		}
	} else {
		params.Format = PictureFormat(conf.DefaultFormat)
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

	log.Debugf("Params: %s\n", params.Dump())

	return params
}
