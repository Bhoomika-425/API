package cache

import (
	"context"
	"encoding/json"
	"project/internal/models"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RDB struct {
	rdb *redis.Client
}
type Cache interface {
	AddingtoCache(ctx context.Context, jid uint, jobdata models.Jobs) error
	GettingCacheData(ctx context.Context, jid uint) (string, error)
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
