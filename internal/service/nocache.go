package service

import (
	"context"
	"minecrafthead/internal/mojang"

	"github.com/google/uuid"
)

// NoCacheSkinStore
type NoCacheSkinStore struct{}

func NewNoCacheSkinStore() *NoCacheSkinStore {
	return &NoCacheSkinStore{}
}

func (s *NoCacheSkinStore) GetHead(ctx context.Context, id uuid.UUID, width, height int, overlay bool) (string, error) {
	return mojang.GetHead64(id, width, height, overlay)
}

// NoCacheUUIDService
type NoCacheUUIDService struct{}

func NewNoCacheUUIDService() *NoCacheUUIDService {
	return &NoCacheUUIDService{}
}

func (s *NoCacheUUIDService) GetUUIDByNickname(ctx context.Context, nick string) (uuid.UUID, error) {
	return mojang.GetUUIDByNickname(nick)
}
