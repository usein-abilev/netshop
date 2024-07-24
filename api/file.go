package api

import (
	"errors"
	"io"
	"log"
	"net/http"
	"netshop/main/db"
	"netshop/main/tools"
	"netshop/main/tools/image"
	"netshop/main/tools/router"
)

var (
	ErrInvalidMultipartFormData = errors.New("invalid multipart/form-data content")
	ErrInvalidFileSize          = errors.New("invalid file size. File must be less than 2MB")
	ErrInvalidImageFormat       = errors.New("invalid image format. File must be in PNG, JPG, JPEG format")
)

// Define the allowed MIME types
var allowedMIMETypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
}

const (
	FilesDirectory = "./static/files/"
	MaxFileSize    = 2 << 20 // 2mb
	ImageQuality   = 80      // 80% quality
)

type fileHandler struct {
	DatabaseConnection *db.DatabaseConnection
	EntityStore        *db.FileEntityStore
}

func InitFileRouter(parent *router.Router, opts *InitEndpointsOptions) {
	handler := fileHandler{
		DatabaseConnection: opts.DatabaseConnection,
		EntityStore:        db.NewFileEntityStore(opts.DatabaseConnection),
	}

	router := parent.Subrouter()

	router.AddRoute("/file/upload", RequireAuth(handler.handleUpload)).
		Methods("POST").
		Name("Upload file").
		Description("This endpoint is used to upload file. It receives file in request body and returns file entity")
}

func (f *fileHandler) handleUpload(w http.ResponseWriter, req *http.Request) {
	err := req.ParseMultipartForm(MaxFileSize)
	if err != nil {
		log.Printf("file/upload: error parsing form: %s", err)
		tools.RespondWithError(w, ErrInvalidFileSize.Error(), http.StatusBadRequest)
		return
	}

	file, _, err := req.FormFile("file")
	if err != nil {
		log.Printf("file/upload: error getting file: %s", err)
		tools.RespondWithError(w, ErrInvalidMultipartFormData.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		log.Printf("file/upload: error reading file: %s", err)
		tools.RespondWithError(w, ErrInvalidImageFormat.Error(), http.StatusInternalServerError)
		return
	}

	contentType := http.DetectContentType(buffer)
	log.Printf("file/upload: Detected content type: %s", contentType)

	if !allowedMIMETypes[contentType] {
		log.Printf("file/upload: invalid file format: %s", contentType)
		tools.RespondWithError(w, ErrInvalidImageFormat.Error(), http.StatusBadRequest)
		return
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		log.Printf("file/upload: error seeking file: %s", err)
		tools.RespondWithError(w, ErrInvalidImageFormat.Error(), http.StatusInternalServerError)
		return
	}

	buffer, err = io.ReadAll(file)
	if err != nil {
		log.Printf("file/upload: error reading file: %s", err)
		tools.RespondWithError(w, ErrInvalidImageFormat.Error(), http.StatusInternalServerError)
		return
	}
	uploadedImage, err := image.UploadImage(buffer, image.UploadImageOptions{
		Dirname:   FilesDirectory,
		Quality:   ImageQuality,
		Watermark: true,
	})
	if err != nil {
		log.Printf("file/upload: error uploading image: %s", err)
		tools.RespondWithError(w, ErrInvalidImageFormat.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("file/upload: File uploaded successfully: %+v", uploadedImage)

	fileEntity, err := f.EntityStore.Create(req.Context(), db.FileEntity{
		Filename:  uploadedImage.Filename,
		Filetype:  uploadedImage.MimeType,
		Path:      uploadedImage.Path,
		Width:     uploadedImage.Width,
		Height:    uploadedImage.Height,
		SizeBytes: uploadedImage.Size,
	})

	if err != nil {
		log.Printf("file/upload: error creating file entity in db: %s", err)
		tools.RespondWithError(w, "Error creating file entity", http.StatusInternalServerError)
		return
	}

	log.Printf("file/upload: File entity created successfully: %+v", fileEntity)
	tools.RespondWithSuccess(w, fileEntity)
}
