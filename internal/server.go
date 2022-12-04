package epubsvc

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func Server(port string, storage *Storage) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/upload", uploadHandler(storage))
	mux.HandleFunc("/download", downloadHandler(storage))

	return http.ListenAndServe(port, mux)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	http.ServeFile(w, r, "internal/index.html")
}

func downloadHandler(s *Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		uuidS := r.URL.Query().Get("uuid")
		path, err := s.Get(uuidS)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fileBytes, err := os.ReadFile(path)
		if err != nil {
			panic(err)
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(fileBytes)
	}
}

func uploadHandler(s *Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		epub, err := storeTempFile(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		u := uuid.New()
		go convertAsync(u, epub, s)
		s.Set(u, epub, false)
		w.Header().Add("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(u.String()))
	}
}
func convertAsync(u uuid.UUID, epub string, s *Storage) {
	root, err := UnpackArchive(epub)
	if err != nil {
		Logger.Error(err.Error())
	}
	path, err := Convert(epub, root)
	if err != nil {
		Logger.Error(err.Error())
	}
	err = DelArchive(root)
	if err != nil {
		Logger.Error(err.Error())
	}
	s.Set(u, path, true)
	Logger.Info("File generated", zap.String("path", epub))
}

func storeTempFile(r *http.Request) (string, error) {
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		return "", err
	}
	defer file.Close()
	Logger.Info("got some shit", zap.String("filename", fileHeader.Filename))

	dst, err := os.Create(fmt.Sprintf("./examples/%s", fileHeader.Filename))
	if err != nil {
		return "", err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return "", err
	}
	return dst.Name(), nil
}
