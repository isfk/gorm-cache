package gormcache

import (
	"context"
	"time"

	"github.com/isfk/go-cache/v3"
	"github.com/isfk/gorm-cache/callbacks"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Cache struct {
	prefix string
	Cache  *cache.Cache
	key    string
	tags   []string
}

type Config struct {
	Prefix    string
	RedisAddr string
	RedisPass string
}

func New(config *Config) *Cache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		Password: config.RedisPass,
	})

	c := cache.New(
		context.Background(),
		rdb,
		cache.WithPrefix(config.Prefix),
		cache.WithExpired(600*time.Second),
	)

	return &Cache{
		prefix: config.Prefix,
		Cache:  c,
	}
}

func (c *Cache) WitchKey(key string) {
	c.key = key
}

func (c *Cache) WithTags(tags ...string) {
	c.tags = append(c.tags, tags...)
}

// Name 实现 gorm.Plugin
func (*Cache) Name() string {
	return "gorm:cache"
}

// Initialize 实现 gorm.Plugin
func (c *Cache) Initialize(db *gorm.DB) (err error) {
	// 注册 callbacks
	err = db.Callback().Create().After("gorm:create").Register("gormcache:after_create", callbacks.AfterCreate(c.Cache))
	if err != nil {
		return err
	}
	err = db.Callback().Update().After("gorm:update").Register("gormcache:after_update", callbacks.AfterUpdate(c.Cache))
	if err != nil {
		return err
	}
	err = db.Callback().Delete().After("gorm:delete").Register("gormcache:after_delete", callbacks.AfterDelete(c.Cache))
	if err != nil {
		return err
	}
	err = db.Callback().Query().Before("gorm:query").Register("gormcache:before_query", callbacks.BeforeQuery(c.Cache))
	if err != nil {
		return err
	}

	// 重写 Query
	err = db.Callback().Query().Replace("gorm:query", callbacks.Query())
	if err != nil {
		return err
	}

	err = db.Callback().Query().After("gorm:query").Register("gormcache:after_query", callbacks.AfterQuery(c.Cache))
	if err != nil {
		return err
	}

	return nil
}
