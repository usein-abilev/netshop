// File package used to manage file conversion, compression and other file related operations
package image

import (
	"crypto/sha1"
	"fmt"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/h2non/bimg"
)

type UploadImageOptions struct {
	// Quality defines the default JPEG quality to be used.
	Quality int
	// Watermark defines the watermark to be used
	Watermark bool
	// Dirname defines the directory name to save the image
	Dirname string
}

type UploadedImage struct {
	Filename string
	Path     string
	MimeType string
	Size     int
	Width    int
	Height   int
}

// UploadImage uploads an image to the specified directory
// This method converts the image to webp format and compresses it
// before saving it to the directory
func UploadImage(buffer []byte, opts UploadImageOptions) (result *UploadedImage, err error) {
	converted, err := convertImageToWebp(buffer, opts.Quality, opts.Watermark)
	if err != nil {
		return nil, err
	}

	imageMetadata, err := bimg.Metadata(converted)
	if err != nil {
		return nil, err
	}

	err = createFolderIfNotExists(opts.Dirname)
	if err != nil {
		return nil, err
	}

	filename := fmt.Sprintf("%s.webp", generateUniqueFilename())
	filePath := path.Join(opts.Dirname, filename)

	err = bimg.Write(filePath, converted)
	if err != nil {
		return nil, err
	}

	result = &UploadedImage{
		Filename: filename,
		Path:     filePath,
		Size:     len(converted),
		MimeType: "image/" + imageMetadata.Type,
		Width:    imageMetadata.Size.Width,
		Height:   imageMetadata.Size.Height,
	}

	return result, nil
}

func generateUniqueFilename() string {
	hash := sha1.New()
	hash.Write([]byte(strconv.FormatInt(time.Now().Unix(), 10)))
	filename := uuid.NewSHA1(uuid.Nil, hash.Sum(nil)).String()
	return filename
}

func createFolderIfNotExists(dirname string) error {
	if _, err := os.Stat(dirname); os.IsNotExist(err) {
		err := os.Mkdir(dirname, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func convertImageToWebp(buffer []byte, quality int, watermark bool) ([]byte, error) {
	options := bimg.Options{
		Quality: quality,
		Type:    bimg.WEBP,
	}

	if watermark {
		options.Watermark = bimg.Watermark{
			Text:    "netshop",
			Opacity: 0.5,
			Width:   100,
			DPI:     72,
		}
	}

	newImageBytes, err := bimg.NewImage(buffer).Process(options)
	if err != nil {
		return nil, err
	}
	return newImageBytes, nil
}
