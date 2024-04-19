package callbacks

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/isfk/go-cache/v3"
	cacheContext "github.com/isfk/gorm-cache/context"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func BeforeQuery(c *cache.Cache) func(tx *gorm.DB) {
	log := slog.With("callback", "before_query")
	return func(tx *gorm.DB) {
		ctx := tx.Statement.Context
		log.Debug("key", "key", ctx.Value(cacheContext.GormCacheKeyCtx{}))

		key := ""
		keyValue := ctx.Value(cacheContext.GormCacheKeyCtx{})
		if keyValue != nil {
			if t, ok := keyValue.(string); ok {
				key = t
			}
		}

		err := c.Get(ctx, key, &tx.Statement.Model)
		if err != nil {
			if errors.Is(err, redis.Nil) {
				return
			}
			log.Error("err", err)
			tx.AddError(err)
			return
		}
		// log.Debug("gorm-cache", slog.Any("dest", tx.Statement.Dest))

		values, err := json.Marshal(tx.Statement.Dest)
		if err != nil {
			tx.AddError(err)
			return
		}

		tx.Statement.Context = context.WithValue(ctx, cacheContext.GormCacheValuesCtx{}, values)
		log.Debug("gorm-cache", "values", string(values))
	}
}
