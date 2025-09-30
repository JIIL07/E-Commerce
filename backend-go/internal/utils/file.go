package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type FileInfo struct {
	Name        string      `json:"name"`
	Size        int64       `json:"size"`
	Path        string      `json:"path"`
	Extension   string      `json:"extension"`
	MimeType    string      `json:"mime_type"`
	MD5Hash     string      `json:"md5_hash"`
	SHA256Hash  string      `json:"sha256_hash"`
	CreatedAt   time.Time   `json:"created_at"`
	ModifiedAt  time.Time   `json:"modified_at"`
	IsDirectory bool        `json:"is_directory"`
	Permissions os.FileMode `json:"permissions"`
}

type FileManager struct {
	BasePath     string
	MaxSize      int64
	AllowedTypes []string
}

func NewFileManager(basePath string, maxSize int64, allowedTypes []string) *FileManager {
	return &FileManager{
		BasePath:     basePath,
		MaxSize:      maxSize,
		AllowedTypes: allowedTypes,
	}
}

func (fm *FileManager) SaveFile(file *multipart.FileHeader, subPath string) (*FileInfo, error) {
	if file.Size > fm.MaxSize {
		return nil, fmt.Errorf("file size exceeds maximum allowed size")
	}

	if !fm.IsAllowedType(file.Filename) {
		return nil, fmt.Errorf("file type not allowed")
	}

	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	fullPath := filepath.Join(fm.BasePath, subPath)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return nil, err
	}

	dst, err := os.Create(fullPath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return nil, err
	}

	return fm.GetFileInfo(fullPath)
}

func (fm *FileManager) GetFileInfo(filePath string) (*FileInfo, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	fileInfo := &FileInfo{
		Name:        info.Name(),
		Size:        info.Size(),
		Path:        filePath,
		Extension:   GetFileExtension(info.Name()),
		MimeType:    GetMimeType(info.Name()),
		CreatedAt:   info.ModTime(),
		ModifiedAt:  info.ModTime(),
		IsDirectory: info.IsDir(),
		Permissions: info.Mode(),
	}

	if !info.IsDir() {
		md5Hash, err := fm.CalculateMD5(filePath)
		if err == nil {
			fileInfo.MD5Hash = md5Hash
		}

		sha256Hash, err := fm.CalculateSHA256(filePath)
		if err == nil {
			fileInfo.SHA256Hash = sha256Hash
		}
	}

	return fileInfo, nil
}

func (fm *FileManager) CalculateMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func (fm *FileManager) CalculateSHA256(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func (fm *FileManager) IsAllowedType(filename string) bool {
	if len(fm.AllowedTypes) == 0 {
		return true
	}

	ext := strings.ToLower(GetFileExtension(filename))
	for _, allowedType := range fm.AllowedTypes {
		if strings.ToLower(allowedType) == ext {
			return true
		}
	}

	return false
}

func (fm *FileManager) DeleteFile(filePath string) error {
	return os.Remove(filePath)
}

func (fm *FileManager) MoveFile(srcPath, dstPath string) error {
	if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
		return err
	}

	return os.Rename(srcPath, dstPath)
}

func (fm *FileManager) CopyFile(srcPath, dstPath string) error {
	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer src.Close()

	if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
		return err
	}

	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}

func (fm *FileManager) ListFiles(dirPath string) ([]*FileInfo, error) {
	var files []*FileInfo

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		fileInfo, err := fm.GetFileInfo(path)
		if err != nil {
			return err
		}

		files = append(files, fileInfo)
		return nil
	})

	return files, err
}

func (fm *FileManager) CreateDirectory(dirPath string) error {
	return os.MkdirAll(dirPath, 0755)
}

func (fm *FileManager) DeleteDirectory(dirPath string) error {
	return os.RemoveAll(dirPath)
}

func (fm *FileManager) GetDirectorySize(dirPath string) (int64, error) {
	var size int64

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			size += info.Size()
		}

		return nil
	})

	return size, err
}

func (fm *FileManager) CleanupOldFiles(dirPath string, olderThan time.Duration) error {
	return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && time.Since(info.ModTime()) > olderThan {
			return os.Remove(path)
		}

		return nil
	})
}

func (fm *FileManager) GetFileContent(filePath string) ([]byte, error) {
	return os.ReadFile(filePath)
}

func (fm *FileManager) WriteFile(filePath string, content []byte) error {
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return err
	}

	return os.WriteFile(filePath, content, 0644)
}

func (fm *FileManager) AppendToFile(filePath string, content []byte) error {
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(content)
	return err
}

func (fm *FileManager) FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

func (fm *FileManager) IsDirectory(filePath string) bool {
	info, err := os.Stat(filePath)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func (fm *FileManager) GetRelativePath(filePath string) string {
	relPath, err := filepath.Rel(fm.BasePath, filePath)
	if err != nil {
		return filePath
	}
	return relPath
}

func (fm *FileManager) GetAbsolutePath(relPath string) string {
	return filepath.Join(fm.BasePath, relPath)
}

func (fm *FileManager) SanitizeFilename(filename string) string {
	filename = strings.ReplaceAll(filename, " ", "_")
	filename = strings.ReplaceAll(filename, "/", "_")
	filename = strings.ReplaceAll(filename, "\\", "_")
	filename = strings.ReplaceAll(filename, ":", "_")
	filename = strings.ReplaceAll(filename, "*", "_")
	filename = strings.ReplaceAll(filename, "?", "_")
	filename = strings.ReplaceAll(filename, "\"", "_")
	filename = strings.ReplaceAll(filename, "<", "_")
	filename = strings.ReplaceAll(filename, ">", "_")
	filename = strings.ReplaceAll(filename, "|", "_")

	return filename
}

func (fm *FileManager) GenerateUniqueFilename(originalName string) string {
	ext := GetFileExtension(originalName)
	baseName := strings.TrimSuffix(originalName, ext)

	timestamp := time.Now().Unix()
	randomStr := GenerateRandomString(8)

	return fmt.Sprintf("%s_%d_%s%s", baseName, timestamp, randomStr, ext)
}

func (fm *FileManager) GetFileStats(dirPath string) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	var totalFiles, totalDirs int64
	var totalSize int64

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			totalDirs++
		} else {
			totalFiles++
			totalSize += info.Size()
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	stats["total_files"] = totalFiles
	stats["total_directories"] = totalDirs
	stats["total_size"] = totalSize
	stats["total_size_formatted"] = FormatBytes(totalSize)

	return stats, nil
}

func (fm *FileManager) BackupFile(filePath string, backupDir string) error {
	if !fm.FileExists(filePath) {
		return fmt.Errorf("file does not exist: %s", filePath)
	}

	fileName := filepath.Base(filePath)
	backupPath := filepath.Join(backupDir, fileName)

	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return err
	}

	return fm.CopyFile(filePath, backupPath)
}

func (fm *FileManager) RestoreFile(backupPath string, restorePath string) error {
	if !fm.FileExists(backupPath) {
		return fmt.Errorf("backup file does not exist: %s", backupPath)
	}

	return fm.CopyFile(backupPath, restorePath)
}

func (fm *FileManager) CompressFile(filePath string, outputPath string) error {
	return fm.CopyFile(filePath, outputPath)
}

func (fm *FileManager) DecompressFile(filePath string, outputPath string) error {
	return fm.CopyFile(filePath, outputPath)
}

func (fm *FileManager) GetFilePermissions(filePath string) (os.FileMode, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}

	return info.Mode(), nil
}

func (fm *FileManager) SetFilePermissions(filePath string, mode os.FileMode) error {
	return os.Chmod(filePath, mode)
}

func (fm *FileManager) GetFileOwner(filePath string) (string, string, error) {
	_, err := os.Stat(filePath)
	if err != nil {
		return "", "", err
	}

	return "", "", nil
}

func (fm *FileManager) SetFileOwner(filePath string, owner, group string) error {
	return nil
}

func (fm *FileManager) GetFileChecksum(filePath string, algorithm string) (string, error) {
	switch strings.ToLower(algorithm) {
	case "md5":
		return fm.CalculateMD5(filePath)
	case "sha256":
		return fm.CalculateSHA256(filePath)
	default:
		return "", fmt.Errorf("unsupported algorithm: %s", algorithm)
	}
}

func (fm *FileManager) VerifyFileChecksum(filePath string, expectedChecksum string, algorithm string) (bool, error) {
	actualChecksum, err := fm.GetFileChecksum(filePath, algorithm)
	if err != nil {
		return false, err
	}

	return actualChecksum == expectedChecksum, nil
}

func (fm *FileManager) GetFileMetadata(filePath string) (map[string]interface{}, error) {
	info, err := fm.GetFileInfo(filePath)
	if err != nil {
		return nil, err
	}

	metadata := map[string]interface{}{
		"name":         info.Name,
		"size":         info.Size,
		"path":         info.Path,
		"extension":    info.Extension,
		"mime_type":    info.MimeType,
		"md5_hash":     info.MD5Hash,
		"sha256_hash":  info.SHA256Hash,
		"created_at":   info.CreatedAt,
		"modified_at":  info.ModifiedAt,
		"is_directory": info.IsDirectory,
		"permissions":  info.Permissions,
	}

	return metadata, nil
}

func (fm *FileManager) SearchFiles(dirPath string, pattern string) ([]*FileInfo, error) {
	var files []*FileInfo

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.Contains(strings.ToLower(info.Name()), strings.ToLower(pattern)) {
			fileInfo, err := fm.GetFileInfo(path)
			if err != nil {
				return err
			}

			files = append(files, fileInfo)
		}

		return nil
	})

	return files, err
}

func (fm *FileManager) GetFileTree(dirPath string) (map[string]interface{}, error) {
	tree := make(map[string]interface{})

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(dirPath, path)
		if err != nil {
			return err
		}

		if relPath == "." {
			return nil
		}

		parts := strings.Split(relPath, string(filepath.Separator))
		current := tree

		for i, part := range parts {
			if i == len(parts)-1 {
				if info.IsDir() {
					current[part] = make(map[string]interface{})
				} else {
					fileInfo, err := fm.GetFileInfo(path)
					if err != nil {
						return err
					}
					current[part] = fileInfo
				}
			} else {
				if _, exists := current[part]; !exists {
					current[part] = make(map[string]interface{})
				}
				current = current[part].(map[string]interface{})
			}
		}

		return nil
	})

	return tree, err
}
