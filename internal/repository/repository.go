package repository

import (
	"context"
	"ogimg/pkg/log"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Repository struct {
	db     *gorm.DB
	rdb    *redis.Client
	logger *log.Logger
}

func NewRepository(logger *log.Logger, db *gorm.DB, conf *viper.Viper) *Repository {
	rdb := redis.NewClient(&redis.Options{
		Addr: conf.GetString("data.redis.addr"),
	})
	return &Repository{
		db:     db,
		rdb:    rdb,
		logger: logger,
	}
}

func (r *Repository) GetFromCache(ctx context.Context, key string) (string, error) {
	r.logger.Info("Get from cache", zap.String("key", key))
	val, err := r.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", err
	}
	return val, nil
}

func (r *Repository) SetToCache(ctx context.Context, key string, val string) error {
	r.logger.Info("Set to cache", zap.String("key", key), zap.String("val", val[:10]))
	err := r.rdb.Set(ctx, key, val, viper.GetDuration("data.redis.expire_time")).Err()
	return err
}

func NewDb() *gorm.DB {
	// TODO: init db
	//db, err := gorm.Open(mysql.Open(conf.GetString("data.mysql.user")), &gorm.Config{})
	//if err != nil {
	//	panic(err)
	//}
	//return db
	return &gorm.DB{}
}
