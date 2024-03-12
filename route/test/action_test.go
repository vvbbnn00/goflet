package test

import (
	"bytes"
	"encoding/json"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/vvbbnn00/goflet/route"
)

var (
	sourcePath   = "/tmp/source.txt"
	targetPath   = "/tmp/target.txt"
	invalidPath  = "/tmp/invalid.txt"
	imagePath    = "/tmp/image.jpg"
	metaFilePath = "/tmp/meta.txt"
	newFilePath  = "/tmp/newfile.txt"
)

// CopyMoveFileRequest is the request body for the copy/move file action
type CopyMoveFileRequest struct {
	OnConflict string `json:"onConflict"`
	SourcePath string `json:"sourcePath"`
	TargetPath string `json:"targetPath"`
}

// CreateFileRequest is the request body for the create file action
type CreateFileRequest struct {
	Path string `json:"path"`
}

func init() {
	router = route.RegisterRoutes()
}

// prepareFileActions prepares the file actions for testing
func TestPrepareFileActions(_ *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/file"+targetPath, nil)
	router.ServeHTTP(w, req)
	req, _ = http.NewRequest(http.MethodDelete, "/file"+invalidPath, nil)
	router.ServeHTTP(w, req)
	req, _ = http.NewRequest(http.MethodDelete, "/file"+newFilePath, nil)
	router.ServeHTTP(w, req)

	postUploadFile(sourcePath, gifData)
	postUploadFile(imagePath, gifData)
	postUploadFile(metaFilePath, gifData)

	time.Sleep(100 * time.Millisecond)
}

func postUploadFile(path string, data []byte) {
	w := httptest.NewRecorder()
	formData := new(bytes.Buffer)
	writer := multipart.NewWriter(formData)
	part, _ := writer.CreateFormFile("file", "test.txt")
	_, err := part.Write(data)
	if err != nil {
		log.Fatalf("Failed to create %s\n", err.Error())
	}
	_ = writer.Close()

	req, _ := http.NewRequest(http.MethodPost, "/file"+path, formData)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	router.ServeHTTP(w, req)
}

func removeFile(path string) {
	req, _ := http.NewRequest(http.MethodDelete, "/api/action/delete", bytes.NewReader([]byte(`{"path":"`+path+`"}`)))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
}

func TestCopyFile(t *testing.T) {
	tests := []struct {
		name           string
		request        CopyMoveFileRequest
		expectedStatus int
	}{
		{"Copy Existing File", CopyMoveFileRequest{OnConflict: "overwrite", SourcePath: sourcePath, TargetPath: targetPath}, http.StatusOK},
		{"Copy Non-Existing Source", CopyMoveFileRequest{OnConflict: "overwrite", SourcePath: invalidPath, TargetPath: targetPath}, http.StatusNotFound},
		{"Copy to Existing Target - Overwrite", CopyMoveFileRequest{OnConflict: "overwrite", SourcePath: sourcePath, TargetPath: targetPath}, http.StatusOK},
		{"Copy to Existing Target - Abort", CopyMoveFileRequest{OnConflict: "abort", SourcePath: sourcePath, TargetPath: targetPath}, http.StatusConflict},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, _ := json.Marshal(tt.request)
			req, _ := http.NewRequest(http.MethodPost, "/api/action/copy", bytes.NewReader(jsonData))
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}

	removeFile(targetPath)
}

func TestCreateFile(t *testing.T) {
	tests := []struct {
		name           string
		request        CreateFileRequest
		expectedStatus int
	}{
		{"Create New File", CreateFileRequest{Path: newFilePath}, http.StatusCreated},
		{"Create Existing File", CreateFileRequest{Path: sourcePath}, http.StatusConflict},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, _ := json.Marshal(tt.request)
			req, _ := http.NewRequest(http.MethodPost, "/api/action/create", bytes.NewReader(jsonData))
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestMoveFile(t *testing.T) {
	tests := []struct {
		name           string
		request        CopyMoveFileRequest
		expectedStatus int
	}{
		{"Move Existing File", CopyMoveFileRequest{OnConflict: "overwrite", SourcePath: sourcePath, TargetPath: targetPath}, http.StatusOK},
		{"Move Non-Existing Source", CopyMoveFileRequest{OnConflict: "overwrite", SourcePath: invalidPath, TargetPath: targetPath}, http.StatusNotFound},
		{"Move to Existing Target - Overwrite", CopyMoveFileRequest{OnConflict: "overwrite", SourcePath: sourcePath, TargetPath: targetPath}, http.StatusOK},
		{"Move to Existing Target - Abort", CopyMoveFileRequest{OnConflict: "abort", SourcePath: sourcePath, TargetPath: targetPath}, http.StatusConflict},
	}

	for _, tt := range tests {
		if tt.name == "Move to Existing Target - Overwrite" || tt.name == "Move to Existing Target - Abort" {
			postUploadFile(sourcePath, gifData)
			time.Sleep(100 * time.Millisecond)
		}
		t.Run(tt.name, func(t *testing.T) {
			jsonData, _ := json.Marshal(tt.request)
			req, _ := http.NewRequest(http.MethodPost, "/api/action/move", bytes.NewReader(jsonData))
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}

	removeFile(targetPath)
}

func TestGetImage(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		query          string
		expectedStatus int
	}{
		{"Get Existing Image", imagePath, "", http.StatusOK},
		{"Get Non-Existing Image", invalidPath, "", http.StatusNotFound},
		{"Get Image with Parameters", imagePath, "?w=100&h=100&q=80&f=jpg&a=90&s=fit", http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, "/api/image/"+tt.path+tt.query, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestGetFileMeta(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		expectedStatus int
	}{
		{"Get Existing File Meta", metaFilePath, http.StatusOK},
		{"Get Non-Existing File Meta", invalidPath, http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, "/api/meta/"+tt.path, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
