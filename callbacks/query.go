package callbacks

import (
	"log/slog"

	"gorm.io/gorm"
)

func Query() func(tx *gorm.DB) {
	log := slog.With("callback", "query")
	return func(tx *gorm.DB) {
		log.Info("query", "RowsAffected", tx.Statement.RowsAffected)
		if tx.Error != nil {
			return
		}
	}
}
