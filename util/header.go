package util

import (
	"errors"
	"github.com/vvbbnn00/goflet/config"
	"strconv"
	"strings"
	"time"
)

// HeaderParseRangeUpload Parse the range header and return the start and end
func HeaderParseRangeUpload(contentRange string, contentLength string) (start int64, end int64, total int64, err error) {
	uploadLimit := config.GofletCfg.FileConfig.UploadLimit

	contentLengthInt, err := strconv.ParseInt(contentLength, 10, 64)
	if err != nil {
		return 0, 0, 0, errors.New("invalid content length")
	}
	// Check if the content length is within the upload limit
	if contentLengthInt > uploadLimit {
		return 0, 0, 0, errors.New("file size exceeds the upload limit")
	}

	// If the range header is empty, return the full content length
	if contentRange == "" {
		return 0, contentLengthInt - 1, contentLengthInt, nil
	}

	// Check if the range header is in the correct format
	if !strings.HasPrefix(contentRange, "bytes ") {
		return 0, 0, 0, errors.New("invalid range header format")
	}

	// Remove the "bytes " prefix
	rangeStr := strings.TrimPrefix(contentRange, "bytes ")

	// Split the range into range and total parts
	rangeTotalParts := strings.Split(rangeStr, "/")
	if len(rangeTotalParts) != 2 {
		return 0, 0, 0, errors.New("invalid range header format")
	}

	// Parse the total part
	total, err = strconv.ParseInt(rangeTotalParts[1], 10, 64)
	if err != nil {
		return 0, 0, 0, errors.New("invalid total value")
	}

	// Split the range into start and end parts
	rangeParts := strings.Split(rangeTotalParts[0], "-")
	if len(rangeParts) != 2 {
		return 0, 0, 0, errors.New("invalid range format")
	}

	// Parse the start part
	if rangeParts[0] != "" {
		start, err = strconv.ParseInt(rangeParts[0], 10, 64)
		if err != nil {
			return 0, 0, 0, errors.New("invalid start value")
		}
	}

	// Parse the end part
	if rangeParts[1] != "" {
		end, err = strconv.ParseInt(rangeParts[1], 10, 64)
		if err != nil {
			return 0, 0, 0, errors.New("invalid end value")
		}
	} else {
		end = start + contentLengthInt - 1
	}

	// Check if the range is valid
	if start > end {
		return 0, 0, 0, errors.New("invalid range: start must be less than or equal to end")
	}

	// Check if the range is within the content length
	if end >= total {
		return 0, 0, 0, errors.New("range exceeds total content length")
	}

	// Check if content length is equal to the upload range
	if contentLengthInt != end-start+1 {
		return 0, 0, 0, errors.New("content length does not match the range")
	}

	// Check if total is within the upload limit
	if total > uploadLimit {
		return 0, 0, 0, errors.New("file size exceeds the upload limit")
	}

	return start, end, total, nil
}

// HeaderParseRangeDownload Parse the range header and return the start and end
func HeaderParseRangeDownload(rangeHeader string, fileSize int64) (start int64, end int64, err error) {
	// If the range header is empty, return the full content length
	if rangeHeader == "" {
		return 0, fileSize - 1, nil
	}

	// Check if the range header is in the correct format
	if !strings.HasPrefix(rangeHeader, "bytes=") {
		return 0, 0, errors.New("invalid range header format")
	}

	// Remove the "bytes=" prefix
	rangeStr := strings.TrimPrefix(rangeHeader, "bytes=")

	// Split the range into start and end parts
	rangeParts := strings.Split(rangeStr, "-")
	if len(rangeParts) != 2 {
		return 0, 0, errors.New("invalid range format")
	}

	// Parse the start part
	if rangeParts[0] != "" {
		start, err = strconv.ParseInt(rangeParts[0], 10, 64)
		if err != nil {
			return 0, 0, errors.New("invalid start value")
		}
	} else {
		// If the start part is empty, parse the end part as the last N bytes
		if rangeParts[1] != "" {
			lastN, err := strconv.ParseInt(rangeParts[1], 10, 64)
			if err != nil {
				return 0, 0, errors.New("invalid end value")
			}
			start = fileSize - lastN
			end = fileSize - 1
			return start, end, nil
		}
	}

	// Parse the end part
	if rangeParts[1] != "" {
		end, err = strconv.ParseInt(rangeParts[1], 10, 64)
		if err != nil {
			return 0, 0, errors.New("invalid end value")
		}
	} else {
		end = fileSize - 1
	}

	// Check if the range is valid
	if start > end {
		return 0, 0, errors.New("invalid range: start must be less than or equal to end")
	}

	// Adjust the end value if it exceeds the file size
	if end >= fileSize {
		end = fileSize - 1
	}

	// Check if the range is within the content length
	if start >= fileSize {
		return 0, 0, errors.New("range exceeds total content length")
	}

	return start, end, nil
}

// HeaderDateToInt64 converts a date string in the format
// "<day-name>, <day> <month> <year> <hour>:<minute>:<second> GMT"
// to an int64 representing the number of seconds since the Unix epoch.
func HeaderDateToInt64(date string) int64 {
	t, err := time.Parse(time.RFC1123, date)
	if err != nil {
		return 0
	}
	return t.Unix()
}

// Int64ToHeaderDate converts an int64 representing the number of seconds
// since the Unix epoch to a date string in the format
// "<day-name>, <day> <month> <year> <hour>:<minute>:<second> GMT".
func Int64ToHeaderDate(timestamp int64) string {
	t := time.Unix(timestamp, 0).UTC()
	// According to https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Last-Modified, the timezone should be GMT
	return strings.Replace(t.Format(time.RFC1123), "UTC", "GMT", 1)
}
