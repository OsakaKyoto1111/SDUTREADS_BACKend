package service

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type FileService struct {
	UploadDir   string
	BaseURL     string
	MaxFileSize int64
	AllowedExts map[string]bool
}

func NewFileService(uploadDir, baseURL string, maxFileSize int64, allowedExts []string) *FileService {
	extMap := make(map[string]bool)
	for _, e := range allowedExts {
		extMap[strings.ToLower(e)] = true
	}
	return &FileService{
		UploadDir:   uploadDir,
		BaseURL:     baseURL,
		MaxFileSize: maxFileSize,
		AllowedExts: extMap,
	}
}

func (fs *FileService) SaveFiles(postID uint, files []*multipart.FileHeader) ([]string, error) {
	uploadDir := filepath.Join(fs.UploadDir, "posts", fmt.Sprint(postID))
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, fmt.Errorf("create upload dir: %w", err)
	}

	var urls []string
	for _, fh := range files {
		if fh.Size > fs.MaxFileSize && fs.MaxFileSize > 0 {
			return nil, fmt.Errorf("file %s too large", fh.Filename)
		}
		ext := strings.ToLower(filepath.Ext(fh.Filename))
		if len(ext) > 0 && ext[0] == '.' {
			ext = ext[1:]
		}
		if len(fs.AllowedExts) > 0 && !fs.AllowedExts[ext] {
			return nil, fmt.Errorf("file %s has forbidden extension", fh.Filename)
		}

		src, err := fh.Open()
		if err != nil {
			return nil, fmt.Errorf("open uploaded file: %w", err)
		}
		defer src.Close()

		newName := uuid.New().String()
		if ext != "" {
			newName = newName + "." + ext
		}
		filePath := filepath.Join(uploadDir, newName)

		dst, err := os.Create(filePath)
		if err != nil {
			return nil, fmt.Errorf("create destination file: %w", err)
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			return nil, fmt.Errorf("copy file content: %w", err)
		}

		relURL := filepath.ToSlash(filepath.Join(fs.BaseURL, "posts", fmt.Sprint(postID), newName))
		urls = append(urls, relURL)
	}

	return urls, nil
}
