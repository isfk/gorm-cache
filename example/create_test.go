package main

import (
	"context"
	"log/slog"
	"testing"

	cache "github.com/isfk/gorm-cache"
	cc "github.com/isfk/gorm-cache/context"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCreate(t *testing.T) {
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

	info := &User{
		Username: "sfk",
	}

	ctx := cc.New(context.Background()).WithKey("keys.user:").WithTags("tags.user:", "tags.user.all").CC()
	err = db.WithContext(ctx).Create(&info).Error
	if err != nil {
		panic(err)
	}
	slog.Debug("create", "info", info)

}
