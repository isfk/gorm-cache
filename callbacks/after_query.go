package callbacks

import (
	"github.com/isfk/go-cache/v3"
	"gorm.io/gorm"
)

func AfterQuery(c *cache.Cache) func(tx *gorm.DB) {
	return func(tx *gorm.DB) {
		ctx := tx.Statement.Context

		key, tags, valueMap, err := GetDataFromCtx(ctx, tx.Statement)
		if err != nil {
			tx.AddError(err)
			return
		}

		if valueMap == nil {
			return
		}

		c.AddTag(tags...)
		err = c.Set(ctx, key, valueMap)
		if err != nil {
			tx.AddError(err)
			return
		}
	}
}
