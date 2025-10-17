package handler

import (
	"fmt"
	"io"
	"koois_core/internal/middleware"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type FileHandler struct{}

const uploadDir = "./upload/tmp"

func (_ *FileHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims, err := middleware.GetClaimsFromContext(r)
	if err != nil {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}
	userId := claims.Sub

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validate MIME
	buffer := make([]byte, 512)
	file.Read(buffer)
	fileType := http.DetectContentType(buffer)
	if fileType != "image/jpeg" && fileType != "image/png" {
		http.Error(w, "invalid mime type", http.StatusInternalServerError)
		return
	}

	filename := uuid.New().String() + filepath.Ext(handler.Filename)
	uploadDir := fmt.Sprintf("%s/%s", uploadDir, userId)

	file.Seek(0, 0)
	os.MkdirAll(uploadDir, os.ModePerm)

	dstPath := filepath.Join(uploadDir, filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		fmt.Print(err)
		http.Error(w, "Error create destination file", http.StatusInternalServerError)
		return

	}
	defer dst.Close()
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Error copy file to destination", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, filename)
}

func cleanAndValidatePath(input string, userId string) (string, error) {
	cleaned := strings.TrimPrefix(input, "/")
	cleaned = filepath.Clean(cleaned)
	absUploadDir, _ := filepath.Abs(uploadDir)
	uploadDir := fmt.Sprintf("%s/%s", uploadDir, userId)
	absTarget, _ := filepath.Abs(filepath.Join(uploadDir, cleaned))
	if !strings.HasPrefix(absTarget, absUploadDir) {
		return "", fmt.Errorf("invalid path")
	}
	return absTarget, nil
}

func (_ *FileHandler) Delete(w http.ResponseWriter, r *http.Request) {
	claims, err := middleware.GetClaimsFromContext(r)
	if err != nil {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}
	userId := claims.Sub

	filename := r.PathValue("id")

	safePath, err := cleanAndValidatePath(filename, userId)
	if err != nil {
		http.Error(w, "Error validate path", http.StatusBadRequest)
		return
	}

	if err := os.Remove(safePath); err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "Error filename not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error delete file", http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File deleted successfully"))
}
