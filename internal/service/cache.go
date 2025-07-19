package service

import (
	"context"
	"log/slog"
	"minecrafthead/internal/mojang"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RedisCacheSkinStore struct {
	redis *redis.Client
}

func NewRedisCacheSkinStore(redis *redis.Client) *RedisCacheSkinStore {
	return &RedisCacheSkinStore{
		redis: redis,
	}
}

func (s *RedisCacheSkinStore) GetHead(ctx context.Context, id uuid.UUID, width, height int, overlay bool) (string, error) {
	slog.Info("getHead", "id", id, "width", width, "height", height, "overlay", overlay)

	res, err := s.redis.Get(ctx, id.String()).Result()
	if err != nil && err != redis.Nil {
		slog.Error("Failed to get head from Redis", "id", id, "err", err)
		return "", err
	}

	if err == redis.Nil {
		slog.Info("getHead miss", "id", id)
		head, err := mojang.GetHead64(id, width, height, overlay)
		if err != nil {
			return "", err
		}

		if err := s.redis.Set(ctx, id.String(), head, time.Hour*24).Err(); err != nil {
			slog.Warn("Failed to cache head in Redis", "id", id, "err", err)
		}

		slog.Info("getHead hit", "id", id)
		return head, nil
	}

	return res, nil
}

// NoCacheUUIDService
type RedisCacheUUIDService struct {
	redis *redis.Client
}

func NewRedisCacheUUIDService(redis *redis.Client) *RedisCacheUUIDService {
	return &RedisCacheUUIDService{
		redis: redis,
	}
}

func (s *RedisCacheUUIDService) GetUUIDByNickname(ctx context.Context, nick string) (uuid.UUID, error) {
	res, err := s.redis.Get(ctx, nick).Result()
	if err != nil && err != redis.Nil {
		slog.Error("Failed to get UUID from Redis", "nick", nick, "err", err)
		return uuid.Nil, err
	}

	if err == redis.Nil {
		slog.Info("getUUID miss", "nick", nick)
		id, err := mojang.GetUUIDByNickname(nick)
		if err != nil {
			return uuid.Nil, err
		}

		if err := s.redis.Set(ctx, nick, id.String(), time.Hour*24).Err(); err != nil {
			slog.Warn("Failed to cache UUID in Redis", "nick", nick, "err", err)
		}

		slog.Info("getUUID hit", "nick", nick)
		return id, nil
	}

	return uuid.Parse(res)
}
