package utils

import (
	"net/http"
	"os"
	"io"
	"github.com/wolfgarnet/logging"
	"time"
	"path/filepath"
	"fmt"
	"mime/multipart"
)

var logger logging.Logger

type UploadResponse struct {
	Name string `json:"name"`
	Filename string
	Rid string `json:"rid"`
	Size int64 `json:"size"`
	Length int64 `json:"length"`
	DeleteUrl string `json:"deleteUtl"`
	ThumbnailUrl string `json:"thumbnailUrl"`
	Url string `json:"url"`
	UploadSession string
}

func HandleUpload(request *http.Request, field string, filenamer func(string, string) string) ([]*UploadResponse, error) {
	logger.Debug("Handling upload from %v", field)

	err := request.ParseMultipartForm(32 << 20)
	if err != nil {
		logger.Warning("UNABLE TO PARSE FORM: %v", err)
		return nil, err
	}

	logger.Debug("FORM: %v", request.MultipartForm)

	files, ok := request.MultipartForm.File[field]
	if !ok {
		logger.Warning("No such field, %v", field)
		return nil, fmt.Errorf("No such field, %v", field)
	}

	/*
	file, header, err := request.FormFile(field)
	if err != nil {
		logger.Warning("'%v': %v", field, err)
		return err
	}
	*/


	sessionName := "default"
	if uploadSession, ok := request.MultipartForm.Value["uploadSession"]; ok {
		if len(uploadSession) > 0 {
			sessionName = uploadSession[0]
		}
	}

	rs := make([]*UploadResponse, 0, len(files))

	for k, file := range files {
		logger.Debug("K======%v", k )
		r, err := getFile(file, filenamer, sessionName)
		if err != nil {
			logger.Warning("Unable to get file, %v", err.Error())
		} else {
			rs = append(rs, r)
		}
	}

	return rs, nil
}

func getFile(header *multipart.FileHeader, filenamer func(string, string) string, sessionName string) (*UploadResponse, error) {

	logger.Debug("HEADER FILE NAME: %v", header.Filename)
	filename := filenamer(header.Filename, sessionName)
	file, err := header.Open()
	if err != nil {
		return nil, err
	}

	defer file.Close()

	os.MkdirAll(filepath.Dir(filename), os.ModePerm)
	out, err := os.Create(filename)
	if err != nil {
		logger.Warning("%v", err)
		return nil, err
	}

	defer out.Close()

	// write the content from POST to the file
	_, err = io.Copy(out, file)
	if err != nil {
		logger.Warning("%v", err)
		return nil, err
	}

	stat, _ := out.Stat()

	r := &UploadResponse{
		Name: header.Filename,
		Filename: filename,
		Size: stat.Size(),
		Length: stat.Size(),
		Url:"/",
		UploadSession: sessionName,
	}

	return r, nil

}
func NewFilenamer(userId string) func(string, string) string {
	return func(base, session string) string{
		t := time.Now()

		year := t.Format("2006")
		month := t.Format("01")

		path := filepath.Join("upload", userId, year, month, session)

		ext := filepath.Ext(base)
		basename := base[0:len(ext)]

		file := filepath.Join(path, base)

		println("FIRST", file)

		cnt := 1
		for FileExists(file) {
			cnt++
			newFilename := fmt.Sprintf("%v_%v%v", basename, cnt, ext)
			file = filepath.Join(path, newFilename)
		}

		return file
	}
}

func FileExists(file string) bool {
	if _, err := os.Stat(file); err == nil {
		return true
	}

	return false
}