package main

import (
	"context"
	"errors"
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

	cache := cache.New(&cache.Config{Prefix: "gorm-cache-example"})
	db.Use(cache)
	db = db.Debug()

	type User struct {
		ID        int64
		Username  string
		CreatedAt int64
		UpdatedAt int64
	}

	db.AutoMigrate(&User{})

	ctx := cc.New(context.Background()).WithKey("keys.user:1").WithTags("tags.user:1", "tags.user.all").CC()
	info := &User{}
	err = db.WithContext(ctx).First(&info, 1).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			panic(err)
		}
	}
	slog.Debug("query", slog.Any("info", info))
}
