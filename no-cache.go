package main

import (
	"context"

	"github.com/google/uuid"
)

type NoCacheSkinStore struct{}

func NewNoCacheSkinStore() *NoCacheSkinStore {
	return &NoCacheSkinStore{}
}

func (s *NoCacheSkinStore) getHead(ctx context.Context, id uuid.UUID, width, height int, overlay bool) (string, error) {
	return getHead64(id, width, height, overlay)
}
