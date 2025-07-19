package service

import (
	"context"

	"github.com/google/uuid"
)

type MinecraftSkinService interface {
	GetHead(ctx context.Context, id uuid.UUID, width, height int, overlay bool) (string, error)
}
