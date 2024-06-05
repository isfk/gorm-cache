package callbacks

import (
	cacheContext "github.com/isfk/gorm-cache/context"
	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
)

func Query() func(tx *gorm.DB) {
	return func(tx *gorm.DB) {
		ctx := tx.Statement.Context
		valsValue := ctx.Value(cacheContext.GormCacheValuesCtx{})
		var values []byte
		if valsValue != nil {
			if t, ok := valsValue.([]byte); ok {
				values = t
			}
		}
		if values != nil {
			tx.Statement.Dest = values
			return
		}
		callbacks.Query(tx)
	}
}
