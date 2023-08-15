package model

import (
	"context"

	"github.com/Goboolean/worker/internal/infrastructure/minio"
)



type ModelAdapter struct {
	storage minio.Storage
}

func New(storage *minio.Storage) *ModelAdapter {
	return &ModelAdapter{
		storage: *storage,
	}
}


func (a *ModelAdapter) NewSession(ctx context.Context, name string) (*ModelSessionImpl, error) {
	f, err := a.storage.GetFile(ctx, name)
	if err != nil {
		return nil, err
	}

	return newSession(f)
}