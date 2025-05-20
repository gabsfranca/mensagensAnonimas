package service

import (
	"errors"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/gabsfranca/mensagensAnonimasRH/models"
	"github.com/gabsfranca/mensagensAnonimasRH/storage"
	"github.com/gin-gonic/gin"
)

type MediaUpload struct {
	Content   string
	Files     []*multipart.FileHeader
	TimeStamp time.Time
}

type AnonymousService struct {
	storage storage.StorageService
}

type MediaFile struct {
	URL  string
	Type models.MediaType
}

func NewAnonymousService(storage storage.StorageService) *AnonymousService {
	return &AnonymousService{storage}
}

func ParseAndValidateForm(c *gin.Context) (*MediaUpload, error) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 50<<20)
	if err := c.Request.ParseMultipartForm(50 << 20); err != nil {
		return nil, errors.New("tamanho total excedeu 50mb")
	}

	content := strings.TrimSpace(c.PostForm("content"))
	if content == "" {
		return nil, errors.New("a mensagem não pode estar vazia")
	}

	if len(content) > 1000 {
		content = content[:1000]
	}

	files := c.Request.MultipartForm.File["files"]

	return &MediaUpload{
		Content:   content,
		Files:     files,
		TimeStamp: time.Now(),
	}, nil
}

func (s *AnonymousService) SaveMediaFiles(files []*multipart.FileHeader) []MediaFile {
	var mediaFiles []MediaFile

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			log.Println("erro ao abrir arquivo: ", err)
			continue
		}
		defer file.Close()

		contentType := fileHeader.Header.Get("Content-Type")
		if !isValidMediaType(contentType) {
			log.Println("tipo de arquivo invalido: ", contentType)
			continue
		}

		fileName, err := s.storage.Save(file, fileHeader.Filename)
		if err != nil {
			log.Println("Erro ao salvar arquivo: ", err)
			continue
		}

		mediaType := strings.Split(contentType, "/")[0]

		var mediaTypeEnum models.MediaType
		switch mediaType {
		case "image":
			mediaTypeEnum = models.Image
		case "video":
			mediaTypeEnum = models.Video
		case "audio":
			mediaTypeEnum = models.Audio
		default:
			mediaTypeEnum = models.Image
		}

		mediaFiles = append(mediaFiles, MediaFile{
			URL:  fileName,
			Type: mediaTypeEnum,
		})

		// mediaURLs = append(mediaURLs, fileName)
	}
	return mediaFiles
}

func isValidMediaType(contentType string) bool {
	allowed := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/gif":  true,
		"video/mp4":  true,
		"audio/mpeg": true,
		"audio/wav":  true,
	}
	return allowed[contentType]
}
