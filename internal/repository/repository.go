package repository

import (
	"context"
	"encoding/json"
	"ogimg/internal/model"
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

func (r *Repository) SetWebsiteOgImgToCache(ctx context.Context, url string, val []byte) error {
	r.logger.Info("Set to cache", zap.String("ogimg:url", url), zap.Int("val_size", len(val)))
	ogImgKey := "ogimg:" + url
	err := r.rdb.Set(ctx, ogImgKey, val, viper.GetDuration("data.redis.expire_time")).Err()
	return err
}

func (r *Repository) GetWebsiteOgImgFromCache(ctx context.Context, url string) ([]byte, error) {
	r.logger.Info("Get from cache", zap.String("ogimg:url", url))
	ogImgKey := "ogimg:" + url
	val, err := r.rdb.Get(ctx, ogImgKey).Bytes()
	return val, err
}

func (r *Repository) SetWebSiteDescToCache(ctx context.Context, url string, val model.WebSiteDescType) error {
	r.logger.Info("Get from cache", zap.String("url", url))
	descKey := "desc:" + url
	jsonVal, err := json.Marshal(val)
	if err != nil {
		return err
	}
	err = r.rdb.Set(ctx, descKey, jsonVal, viper.GetDuration("data.redis.expire_time")).Err()
	return err
}

func (r *Repository) GetWebSiteDescToCache(ctx context.Context, url string) (model.WebSiteDescType, error) {
	r.logger.Info("Get from cache", zap.String("desc:url", url))
	desKey := "desc:" + url
	val, err := r.rdb.Get(ctx, desKey).Result()
	if err == redis.Nil {
		return model.WebSiteDescType{}, nil
	} else if err != nil {
		return model.WebSiteDescType{}, err
	}
	var desc model.WebSiteDescType
	err = json.Unmarshal([]byte(val), &desc)
	if err != nil {
		return model.WebSiteDescType{}, err
	}
	return desc, nil
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
