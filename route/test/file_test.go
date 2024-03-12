package test

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/vvbbnn00/goflet/config"
	"github.com/vvbbnn00/goflet/route"
	"github.com/vvbbnn00/goflet/util"
)

var (
	tmpFileName string
	router      *gin.Engine
	gifData     = []byte("GIF89a\x01\x00\x01\x00\x80\x00\x00\xff\xff\xff\x00\x00\x00!\xf9\x04\x01\x00\x00\x00\x00,\x00\x00\x00\x00\x01\x00\x01\x00\x00\x02\x02D\x01\x00;")
	gifLen      = len(gifData)
)

func init() {
	config.InitConfig()
	*config.GofletCfg.JWTConfig.Enabled = false
	router = route.RegisterRoutes()
	prepareFileUpload()
}

func prepareFileUpload() {
	tmpFileName = "/tmp/" + util.RandomString(16) + ".txt"
	req, _ := http.NewRequest(http.MethodDelete, "/file"+tmpFileName, nil)
	router.ServeHTTP(httptest.NewRecorder(), req)
}

// TestFileUpload tests the file upload functionality
func TestFileUpload(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		endpoint       string
		body           []byte
		contentRange   string
		contentLength  string
		expectedStatus int
		expectedBody   string
	}{
		{"Partial Upload Part 1", http.MethodPut, "/upload" + tmpFileName, gifData[:10], "bytes 0-9/" + strconv.Itoa(gifLen), "10", http.StatusAccepted, ""},
		{"Partial Upload Part 2", http.MethodPut, "/upload" + tmpFileName, gifData[10:], "bytes 10-" + strconv.Itoa(gifLen-1) + "/" + strconv.Itoa(gifLen), strconv.Itoa(gifLen - 10), http.StatusAccepted, ""},
		{"Complete Upload", http.MethodPost, "/upload" + tmpFileName, nil, "", "", http.StatusCreated, ""},
		{"Get File", http.MethodGet, "/file" + tmpFileName, nil, "", "", http.StatusOK, string(gifData)},
		{"Get File Part", http.MethodGet, "/file" + tmpFileName, nil, "bytes=0-9", "10", http.StatusPartialContent, string(gifData[:10])},
		{"Delete File", http.MethodDelete, "/file" + tmpFileName, nil, "", "", http.StatusNoContent, ""},
		{"Post File Upload", http.MethodPost, "/file" + tmpFileName, gifData, "", "", http.StatusCreated, ""},
		{"Get File", http.MethodGet, "/file" + tmpFileName, nil, "", "", http.StatusOK, string(gifData)},
		{"Delete File", http.MethodDelete, "/file" + tmpFileName, nil, "", "", http.StatusNoContent, ""},
	}

	for _, tt := range tests {
		if tt.name == "Post File Upload" {
			testFilePostUpload(t)
			time.Sleep(100 * time.Millisecond) // Wait for the file to be written
			continue
		}
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			if tt.body != nil {
				req, _ = http.NewRequest(tt.method, tt.endpoint, bytes.NewReader(tt.body))
			} else {
				req, _ = http.NewRequest(tt.method, tt.endpoint, nil)
			}
			if tt.contentRange != "" {
				if tt.method == http.MethodPut {
					req.Header.Set("Content-Range", tt.contentRange)
				} else {
					req.Header.Set("Range", tt.contentRange)
				}
			}
			if tt.contentLength != "" {
				req.Header.Set("Content-Length", tt.contentLength)
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, tt.expectedBody, w.Body.String())

			if tt.name == "Complete Upload" {
				time.Sleep(100 * time.Millisecond) // Wait for the file to be written
			}
		})
	}
}

func testFilePostUpload(t *testing.T) {
	w := httptest.NewRecorder()
	formData := new(bytes.Buffer)
	writer := multipart.NewWriter(formData)
	part, _ := writer.CreateFormFile("file", "test.txt")
	write, err := part.Write(gifData)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, gifLen, write)
	_ = writer.Close()

	req, _ := http.NewRequest(http.MethodPost, "/file"+tmpFileName, formData)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "", w.Body.String())
}
