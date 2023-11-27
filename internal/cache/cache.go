package cache

import (
	"context"
	"encoding/json"
	"errors"
	"project/internal/models"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type RDB struct {
	rdb *redis.Client
}

//go:generate mockgen -source=cache.go -destination=cache_mock.go -package=cache
type Cache interface {
	AddingtoCache(ctx context.Context, jid uint, jobdata models.Jobs) error
	GettingCacheData(ctx context.Context, jid uint) (string, error)
	Saveotptoredis(ctx context.Context, email string, otp string) (string, error)
	Getaotp(ctx context.Context, email string) (string, error)
	Deleteotp(ctx context.Context, email string) error
}

func NewRDBLayer(rdb *redis.Client) Cache {
	return &RDB{
		rdb: rdb,
	}
}

func (r RDB) AddingtoCache(ctx context.Context, jid uint, jobdata models.Jobs) error {
	jobId := strconv.FormatUint(uint64(jid), 10)
	val, err := json.Marshal(jobdata)
	if err != nil {
		return err
	}
	err = r.rdb.Set(ctx, jobId, val, 2*time.Minute).Err()
	return err
}
func (r RDB) GettingCacheData(ctx context.Context, jid uint) (string, error) {
	jobId := strconv.FormatUint(uint64(jid), 10)
	str, err := r.rdb.Get(ctx, jobId).Result()
	return str, err
}
func (r RDB) Saveotptoredis(ctx context.Context, email string, otp string) (string, error) {
	err := r.rdb.Set(ctx, email, otp, time.Minute*5).Err()
	if err != nil {
		log.Error().Err(err).Str("email", email).Msg("error while saving the otp to Redis")
		return "", err
	}
	return otp, nil
}
func (r RDB) Getaotp(ctx context.Context, email string) (string, error) {
	s, err := r.rdb.Get(ctx, email).Result()
	return s, err
}

func (r RDB) Deleteotp(ctx context.Context, email string) error {
	err := r.rdb.Del(ctx, email).Err()
	if err != nil {
		return errors.New("error in deleting otp")
	}
	return nil
}
