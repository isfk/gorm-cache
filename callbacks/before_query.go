package callbacks

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/isfk/go-cache/v3"
	cacheContext "github.com/isfk/gorm-cache/context"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func BeforeQuery(c *cache.Cache) func(tx *gorm.DB) {
	return func(tx *gorm.DB) {
		ctx := tx.Statement.Context

		key, _, _, err := GetDataFromCtx(ctx, tx.Statement)
		if err != nil {
			tx.AddError(err)
			return
		}
		err = c.Get(ctx, key, tx.Statement.Dest)
		if err != nil {
			if !errors.Is(err, redis.Nil) {
				tx.AddError(err)
			}
			return
		}

		values, err := json.Marshal(tx.Statement.Dest)
		if err != nil {
			tx.AddError(err)
			return
		}

		tx.Statement.Context = context.WithValue(ctx, cacheContext.GormCacheValuesCtx{}, values)
	}
}
