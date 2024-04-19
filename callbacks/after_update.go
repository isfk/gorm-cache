package callbacks

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/isfk/go-cache/v3"
	cacheContext "github.com/isfk/gorm-cache/context"
	"gorm.io/gorm"
)

func AfterUpdate(c *cache.Cache) func(tx *gorm.DB) {
	log := slog.With("callback", "after_update")
	return func(tx *gorm.DB) {
		ctx := tx.Statement.Context

		key := ""
		keyValue := ctx.Value(cacheContext.GormCacheKeyCtx{})
		if keyValue != nil {
			if t, ok := keyValue.(string); ok {
				key = t
			}
		}

		tags := []string{}
		tagsValue := ctx.Value(cacheContext.GormCacheTagsCtx{})
		if tagsValue != nil {
			if t, ok := tagsValue.([]string); ok {
				tags = t
			}
		}

		if tx.Error != nil {
			return
		}

		log.Debug("gorm-cache", slog.Any("dest", tx.Statement.Dest))

		values, err := json.Marshal(tx.Statement.Dest)
		if err != nil {
			tx.AddError(err)
			return
		}

		valueMap := map[string]interface{}{}
		err = json.Unmarshal(values, &valueMap)
		if err != nil {
			tx.AddError(err)
			return
		}

		if strings.HasSuffix(key, ":") {
			switch valueMap[tx.Statement.Schema.PrioritizedPrimaryField.Name].(type) {
			case float64:
				key = fmt.Sprintf("%s%d", key, int64(valueMap[tx.Statement.Schema.PrioritizedPrimaryField.Name].(float64)))
			case int64:
				key = fmt.Sprintf("%s%d", key, valueMap[tx.Statement.Schema.PrioritizedPrimaryField.Name].(int64))
			case string:
				key = fmt.Sprintf("%s%s", key, valueMap[tx.Statement.Schema.PrioritizedPrimaryField.Name].(string))
			}
		}

		for i, tag := range tags {
			if strings.HasSuffix(tag, ":") {
				switch valueMap[tx.Statement.Schema.PrioritizedPrimaryField.Name].(type) {
				case float64:
					tags[i] = fmt.Sprintf("%s%d", tag, int64(valueMap[tx.Statement.Schema.PrioritizedPrimaryField.Name].(float64)))
				case int64:
					tags[i] = fmt.Sprintf("%s%d", tag, valueMap[tx.Statement.Schema.PrioritizedPrimaryField.Name].(int64))
				case string:
					tags[i] = fmt.Sprintf("%s%s", tag, valueMap[tx.Statement.Schema.PrioritizedPrimaryField.Name].(string))
				}
			}
		}

		c.AddTag(tags...)
		err = c.Set(ctx, key, tx.Statement.Dest)
		if err != nil {
			tx.AddError(err)
			return
		}

		log.Debug("gorm-cache", slog.String("done", "."))
	}
}
