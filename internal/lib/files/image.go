package files

import (
	"fmt"
	"github.com/Verce11o/yata/internal/lib/response"
	"github.com/google/uuid"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

func PrepareImage(image *multipart.FileHeader) (string, []byte, string, error) {
	file, err := image.Open()
	defer file.Close()

	if err != nil {
		return "", nil, "", err
	}

	bytes, err := io.ReadAll(file)

	if err != nil {
		return "", nil, "", err
	}

	contentType := http.DetectContentType(bytes)

	if !checkImageMime(contentType) {
		return "", nil, "", response.ErrInvalidImage
	}

	return contentType, bytes, generateImageName(getExtension(contentType)), nil
}

func checkImageMime(imageMime string) bool {
	var imageMimeTypes = map[string]struct{}{
		"image/gif":  {},
		"image/jpeg": {},
		"image/png":  {},
		"image/webp": {},
	}

	_, ok := imageMimeTypes[imageMime]
	return ok
}

func getExtension(contentType string) string {
	return strings.Split(contentType, "/")[1]
}

func generateImageName(extension string) string {
	return fmt.Sprintf("%v.%v", uuid.NewString(), extension)
}
