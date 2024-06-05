package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"testing"

	cache "github.com/isfk/gorm-cache"
	cc "github.com/isfk/gorm-cache/context"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestQuery(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	db, err := gorm.Open(sqlite.Open("./example.db"), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})
	if err != nil {
		panic(err)
	}

	cache := cache.New(&cache.Config{
		Prefix:    "gorm-cache-example",
		RedisAddr: "127.0.0.1:6379",
		RedisPass: "1234567890",
	})
	db.Use(cache)
	db = db.Debug()

	type User struct {
		ID        int64
		Username  string
		CreatedAt int64
		UpdatedAt int64
	}

	db.AutoMigrate(&User{})

	id := 193

	ctx := cc.New(context.Background()).WithKey(fmt.Sprintf("keys.user:%d", id)).WithTags(fmt.Sprintf("tags.user:%d", id), "tags.user.all").CC()
	info := &User{}
	err = db.WithContext(ctx).First(&info, id).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			panic(err)
		}
	}
	slog.Debug("query", slog.Any("info", info))
}
