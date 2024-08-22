package repository

import (
	"context"

	"github.com/SinnerUfa/practicum-metric/internal/metrics"
	"github.com/SinnerUfa/practicum-metric/internal/repository/database"
	"github.com/SinnerUfa/practicum-metric/internal/repository/memory"
	"github.com/SinnerUfa/practicum-metric/internal/repository/unload"
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

type Config struct {
	StoreInterval   uint
	FileStoragePath string
	Restore         bool
	DatabaseDSN     string
}

type Repository struct {
	storageType RepositoryType
	storage     any
}

func New(ctx context.Context, cfg Config) (*Repository, error) {
	r := &Repository{}
	if cfg.DatabaseDSN != "" {
		r.storageType = DBStorageType
		storage, err := database.New(ctx, cfg.DatabaseDSN)
		if err != nil {
			return nil, err
		}
		r.storage = storage
		return r, nil
	}
	if cfg.Restore {
		r.storageType = UnloadStorageType
		storage, err := unload.New(ctx, cfg.FileStoragePath, cfg.StoreInterval)
		if err != nil {
			return nil, err
		}
		r.storage = storage
		return r, nil
	}
	r.storageType = MemoryStorageType
	r.storage = memory.New()
	return r, nil
}

type MemoryStorageAdapter struct {
	*memory.Memory
}

func (ma *MemoryStorageAdapter) Set(_ context.Context, m metrics.Metric) error {
	return ma.Memory.Set(m)
}

func (ma *MemoryStorageAdapter) Get(_ context.Context, m *metrics.Metric) error {
	return ma.Memory.Get(m)
}

func (ma *MemoryStorageAdapter) SetList(_ context.Context, in []metrics.Metric) error {
	return ma.Memory.SetList(in)
}

func (ma *MemoryStorageAdapter) GetList(_ context.Context) ([]metrics.Metric, error) {
	return ma.Memory.GetList()
}

type UnloadStorageAdapter struct {
	*unload.Unload
}

func (ua *UnloadStorageAdapter) Set(_ context.Context, m metrics.Metric) error {
	return ua.Unload.Set(m)
}

func (ua *UnloadStorageAdapter) Get(_ context.Context, m *metrics.Metric) error {
	return ua.Unload.Get(m)
}

func (ua *UnloadStorageAdapter) SetList(_ context.Context, in []metrics.Metric) error {
	return ua.Unload.SetList(in)
}

func (ua *UnloadStorageAdapter) GetList(_ context.Context) ([]metrics.Metric, error) {
	return ua.Unload.GetList()
}

type DBStorageAdapter struct {
	*database.Database
}

func (dba *DBStorageAdapter) Set(ctx context.Context, m metrics.Metric) error {
	return dba.Database.Set(ctx, m)
}

func (dba *DBStorageAdapter) Get(ctx context.Context, m *metrics.Metric) error {
	return dba.Database.Get(ctx, m)
}

func (dba *DBStorageAdapter) SetList(ctx context.Context, in []metrics.Metric) error {
	return dba.Database.SetList(ctx, in)
}

func (dba *DBStorageAdapter) GetList(ctx context.Context) (out []metrics.Metric, err error) {
	return dba.Database.GetList(ctx)
}

func (r *Repository) Storage() Storage {
	switch r.storageType {
	case DBStorageType:
		return &DBStorageAdapter{Database: r.storage.(*(database.Database))}
	case UnloadStorageType:
		return &UnloadStorageAdapter{Unload: r.storage.(*(unload.Unload))}
	default:
		return &MemoryStorageAdapter{Memory: r.storage.(*(memory.Memory))}
	}
}

func (r *Repository) Close() error {
	switch r.storageType {
	case MemoryStorageType:
		return nil
	case DBStorageType:
		return r.storage.(*(database.Database)).Close()
	case UnloadStorageType:
		return r.storage.(*(unload.Unload)).Close()
	}
	return nil
}

func (r *Repository) Type() string {
	switch r.storageType {
	case MemoryStorageType:
		return "memory storage"
	case DBStorageType:
		return "DB storage"
	case UnloadStorageType:
		return "file storage"
	}
	return "unknown storage"
}
