package callbacks

import (
	"log/slog"

	cacheContext "github.com/isfk/gorm-cache/context"
	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
)

func Query() func(tx *gorm.DB) {
	log := slog.With("callback", "query")
	return func(tx *gorm.DB) {
		ctx := tx.Statement.Context
		valsValue := ctx.Value(cacheContext.GormCacheValuesCtx{})
		var values []byte
		if valsValue != nil {
			if t, ok := valsValue.([]byte); ok {
				values = t
			}
		}
		log.Debug("gorm-cache", "values", values)
		if len(string(values)) > 0 {
			tx.Statement.Dest = values
			return
		}

		log.Debug("gorm-cache", "msg", "nocache")
		callbacks.Query(tx)
	}
}
