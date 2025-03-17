package storage

import (
	"context"
	"encoding/json"
	"github.com/pasha1980/asanaclient/config"
	"github.com/pasha1980/asanaclient/internal/service"
	"os"
	"sync"
)

type storage struct {
	basePath string
	muMap    map[string]*sync.Mutex
}

func (s *storage) Save(ctx context.Context, key string, data any) error {
	mu, ok := s.muMap[key]
	if !ok {
		mu = &sync.Mutex{}
		s.muMap[key] = mu
	}

	mu.Lock()
	defer mu.Unlock()

	f, err := os.OpenFile(s.basePath+"/"+key+".json", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	return json.NewEncoder(f).Encode(data)
}

func NewStorage() service.Storage {
	cfg := config.Get()

	return &storage{
		basePath: cfg.StorageBasePath,
		muMap:    make(map[string]*sync.Mutex),
	}
}
