package photos

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	testPhotoFileName = "photo.png"
	testPhotoBase64   = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8z8BQDwAEhQGAhKmMIQAAAABJRU5ErkJggg=="
)

// CreateTempPNG создаёт временный PNG-файл для загрузки в объявление и возвращает путь к нему
func CreateTempPNG(t *testing.T) string {
	t.Helper()

	photoBytes, err := base64.StdEncoding.DecodeString(testPhotoBase64)
	require.NoError(t, err, "failed to decode test photo")

	filePath := filepath.Join(t.TempDir(), testPhotoFileName)
	require.NoError(t, os.WriteFile(filePath, photoBytes, 0o600), "failed to create test photo")

	return filePath
}

// CreateTempPNGs создаёт указанное количество временных PNG-файлов и возвращает пути к ним
func CreateTempPNGs(t *testing.T, count int) []string {
	t.Helper()

	photoPaths := make([]string, 0, count)
	for range count {
		photoPaths = append(photoPaths, CreateTempPNG(t))
	}

	return photoPaths
}
