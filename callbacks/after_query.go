package callbacks

import (
	"log/slog"

	cacheContext "github.com/isfk/gorm-cache/context"
	"gorm.io/gorm"
)

func AfterQuery() func(tx *gorm.DB) {
	log := slog.With("callback", "after_query")
	return func(tx *gorm.DB) {
		if tx.Error != nil {
			return
		}
		// 存储缓存
		ctx := tx.Statement.Context
		log.Debug("key", ctx.Value(cacheContext.GormCacheKeyCtx{}))

		if tx.Error != nil {
			return
		}
	}
}
