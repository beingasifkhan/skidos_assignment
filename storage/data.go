package storage

import (
	"io"
	"mime/multipart"
	"os"
)

type Storage struct {
	StoragePath string
}

func NewStorage(storagePath string) *Storage {
	return &Storage{StoragePath: storagePath}
}

func (s *Storage) SaveVideo(fileHeader *multipart.FileHeader) error {
	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()
	dstFile, err := os.Create(s.StoragePath + fileHeader.Filename)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	_, err = io.Copy(dstFile, file)
	if err != nil {
		return err
	}

	return nil
}
