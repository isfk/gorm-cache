package callbacks

import (
	"log"

	cacheContext "github.com/isfk/gorm-cache/context"
	"gorm.io/gorm"
)

func AfterDelete() func(tx *gorm.DB) {
	return func(tx *gorm.DB) {
		log.Println("table: ", tx.Statement.Table)
		log.Println("table: ", tx.Statement.Schema)
		log.Println("table: ", tx.Statement.Schema.Table)

		ctx := tx.Statement.Context
		log.Println("key", ctx.Value(cacheContext.GormCacheKeyCtx{}))
		log.Println("tags", ctx.Value(cacheContext.GormCacheTagsCtx{}))

		if tx.Error != nil {
			return
		}

		// 清除旧缓存
	}
}
