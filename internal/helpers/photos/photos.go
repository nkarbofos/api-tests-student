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
	require.NoError(t, err, "не удалось декодировать тестовую фотографию")

	filePath := filepath.Join(t.TempDir(), testPhotoFileName)
	require.NoError(t, os.WriteFile(filePath, photoBytes, 0o600), "не удалось создать тестовую фотографию")

	return filePath
}

// CreateTempPNGs создаёт указанное количество временных PNG-файлов и возвращает пути к ним
func CreateTempPNGs(t *testing.T, count int) []string {
	t.Helper()

	photoPaths := make([]string, 0, count)
	for i := 0; i < count; i++ {
		photoPaths = append(photoPaths, CreateTempPNG(t))
	}

	return photoPaths
}
