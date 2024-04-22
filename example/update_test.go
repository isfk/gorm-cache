package main

import (
	"context"
	"fmt"
	"log/slog"
	"testing"

	cache "github.com/isfk/gorm-cache"
	cc "github.com/isfk/gorm-cache/context"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestUpdate(t *testing.T) {
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
		ID:       113,
		Username: "sfk",
	}

	info.Username = fmt.Sprintf("user-update-%d", 113)
	ctx := cc.New(context.Background()).WithKey(fmt.Sprintf("keys.user:%d", 113)).WithTags(fmt.Sprintf("tags.user:%d", 113), "tags.user.all").CC()
	err = db.WithContext(ctx).Save(&info).Error
	if err != nil {
		panic(err)
	}
	slog.Debug("update", "info", info)

}
