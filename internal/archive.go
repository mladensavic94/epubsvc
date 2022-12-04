package epubsvc

import (
	"archive/zip"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
)

func UnpackArchive(epub string) (string, error) {
	dst := strings.Replace(epub, ".epub", "", -1)
	archive, err := zip.OpenReader(epub)
	if err != nil {
		panic(err)
	}
	for _, f := range archive.File {
		filePath := filepath.Join(dst, f.Name)
		Logger.Info("unzipping file", zap.String("file", filePath))

		if !strings.HasPrefix(filePath, filepath.Clean(dst)+string(os.PathSeparator)) {
			Logger.Info("invalid file path", zap.String("file", filePath))
			return "", errors.New("unzip err")
		}
		if f.FileInfo().IsDir() {
			Logger.Info("creating directory...", zap.String("file", filePath))
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return "", errors.New("unzip err")
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return "", errors.New("unzip err " + err.Error())
		}

		fileInArchive, err := f.Open()
		if err != nil {
			return "", errors.New("unzip err " + err.Error())
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			return "", errors.New("unzip err " + err.Error())
		}

		dstFile.Close()
		fileInArchive.Close()
	}
	return dst, nil
}

func DelArchive(path string) error {
	return os.RemoveAll(path)
}
