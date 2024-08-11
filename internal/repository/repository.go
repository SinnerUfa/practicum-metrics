package repository

import (
	"context"

	"github.com/SinnerUfa/practicum-metric/internal/repository/memory"

	"github.com/SinnerUfa/practicum-metric/internal/repository/unload"

	"github.com/SinnerUfa/practicum-metric/internal/metrics"
)

type RepositoryType int

const (
	MemoryStorageType RepositoryType = iota
	UnloadStorageType
	DBStorageType
)

type Storage interface {
	metrics.Setter
	metrics.Getter
	metrics.ListSetter
	metrics.ListGetter
}

type Repository struct {
	storageType RepositoryType
	storage     any
}

type Config struct {
	StoreInterval   uint
	FileStoragePath string
	Restore         bool
	DatabaseDSN     string
}

func New(ctx context.Context, cfg Config) (*Repository, error) {
	r := &Repository{}
	if cfg.Restore {
		r.storageType = UnloadStorageType
		storage, err := unload.New(ctx, cfg.FileStoragePath, cfg.StoreInterval)
		if err != nil {
			return nil, err
		}
		r.storage = storage
	}
	r.storageType = MemoryStorageType
	r.storage = memory.New()
	return r, nil
}

func (r *Repository) Storage() Storage {
	return r.storage.(Storage)
}

func (r *Repository) Close() error {
	switch r.storageType {
	case MemoryStorageType, DBStorageType:
		return nil
	case UnloadStorageType:
		return r.storage.(*(unload.Unload)).Close()
	}
	return nil
}
