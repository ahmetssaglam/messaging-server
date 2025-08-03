package database

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"messaging-server/internal/configs"
	log "messaging-server/internal/logging"
	"messaging-server/internal/models"
	"time"
)

type RedisClientTemplate struct {
	client *redis.Client
	ctx    context.Context
	ttl    time.Duration
}

var RedisClient *RedisClientTemplate

// ConnectRedis initializes the singleton RedisClient
func ConnectRedis() error {
	if RedisClient != nil {
		return nil
	}

	cfg := configs.RedisConfig

	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	rdb := redis.NewClient(&redis.Options{
		Addr:         addr,
		DB:           cfg.DB,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	// fail-fast ping
	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis ping failed: %w", err)
	}

	RedisClient = &RedisClientTemplate{client: rdb, ctx: ctx, ttl: time.Duration(cfg.TTL) * time.Second}
	log.Logger.Info("Redis connection established")
	return nil
}

// ensureConnection checks if the Redis connection is alive and reconnects if not
func (r *RedisClientTemplate) ensureConnection() {
	if err := r.client.Ping(r.ctx).Err(); err != nil {
		log.Logger.Warningf("lost Redis connection (%v), reconnecting…", err)
		_ = r.client.Close()

		// clear so ConnectRedis will re-init
		RedisClient = nil
		if err := ConnectRedis(); err != nil {
			log.Logger.Fatalf("Redis reconnect failed: %v", err)
		}
		// reset receiver to the fresh global
		*r = *RedisClient
	}
}

// InsertRecord inserts a new record into Redis with the specified TTL
func (r *RedisClientTemplate) InsertRecord(rec models.RedisRecord) error {
	r.ensureConnection()
	if err := r.client.Set(r.ctx, rec.MessageID, rec.SentAt, r.ttl).Err(); err != nil {
		return fmt.Errorf("redis SET failed: %w", err)
	}
	log.Logger.Debugf("Redis SET succeeded for key: %s", rec.MessageID)
	return nil
}
