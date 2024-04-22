package callbacks

import (
	"log/slog"

	cacheContext "github.com/isfk/gorm-cache/context"
	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
)

func Query() func(tx *gorm.DB) {
	slog.Debug("query", "start", ".")
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
		slog.Debug("query", "msg", "nocache")
		callbacks.Query(tx)
		slog.Debug("query", "done", ".")
	}
}
