package service

import (
	"context"
	"minecrafthead/internal/mojang"

	"github.com/google/uuid"
)

type NoCacheSkinStore struct{}

func NewNoCacheSkinStore() *NoCacheSkinStore {
	return &NoCacheSkinStore{}
}

func (s *NoCacheSkinStore) GetHead(ctx context.Context, id uuid.UUID, width, height int, overlay bool) (string, error) {
	return mojang.GetHead64(id, width, height, overlay)
}
