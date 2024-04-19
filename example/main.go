package main

import (
	"context"
	"fmt"
	"log/slog"

	cache "github.com/isfk/gorm-cache"
	cc "github.com/isfk/gorm-cache/context"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
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

	id := info.ID
	ctx = cc.New(context.Background()).WithKey(fmt.Sprintf("keys.user:%d", id)).WithTags(fmt.Sprintf("tags.user:%d", id), "tags.user.all").CC()
	info = &User{}
	err = db.WithContext(ctx).First(&info, id).Error
	if err != nil {
		panic(err)
	}
	slog.Debug("query", "info", info)

	info.Username = fmt.Sprintf("user-%d", id)
	ctx = cc.New(context.Background()).WithKey(fmt.Sprintf("keys.user:%d", id)).WithTags(fmt.Sprintf("tags.user:%d", id), "tags.user.all").CC()
	err = db.WithContext(ctx).Save(&info).Error
	if err != nil {
		panic(err)
	}
	slog.Debug("update", "info", info)

	ctx = cc.New(context.Background()).WithKey("keys.user:1").WithTags("tags.user:1", "tags.user.all").CC()
	err = db.WithContext(ctx).Delete(&User{}, 1).Error
	if err != nil {
		panic(err)
	}
	slog.Debug("delete")
}
