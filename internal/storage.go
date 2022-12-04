package epubsvc

import (
	"errors"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Storage struct {
	cache map[uuid.UUID]FileWrap
	path  map[string]uuid.UUID
}

func NewStorage() *Storage {
	return &Storage{
		cache: make(map[uuid.UUID]FileWrap),
		path:  make(map[string]uuid.UUID),
	}
}

type FileWrap struct {
	filepath string
	ready    bool
}

func (s *Storage) Get(uuidS string) (string, error) {
	u, err := uuid.Parse(uuidS)
	if err != nil {
		Logger.Error("Cant parse UUID", zap.String("UUID", uuidS))
		return "", err
	}
	fw, found := s.cache[u]
	if !found {
		Logger.Error("Cant find file with that UUID", zap.String("UUID", uuidS))
		return "", errors.New("cant find file with that UUID " + uuidS)
	}
	if !fw.ready {
		Logger.Error("File not ready", zap.String("UUID", uuidS))
		return "", errors.New("File not ready with that UUID " + uuidS)
	}
	return fw.filepath, nil
}

func (s *Storage) Set(u uuid.UUID, path string, ready bool) {
	s.cache[u] = FileWrap{path, ready}
}
