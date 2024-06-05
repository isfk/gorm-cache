package callbacks

import (
	"github.com/isfk/go-cache/v3"
	"gorm.io/gorm"
)

func AfterDelete(c *cache.Cache) func(tx *gorm.DB) {
	return func(tx *gorm.DB) {
		ctx := tx.Statement.Context

		key, tags, _, err := GetDataFromCtx(ctx, tx.Statement)
		if err != nil {
			tx.AddError(err)
			return
		}

		c.AddTag(tags...)
		err = c.Del(ctx, key)
		if err != nil {
			tx.AddError(err)
			return
		}
	}
}
