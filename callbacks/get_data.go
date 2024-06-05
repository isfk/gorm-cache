package callbacks

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	cacheContext "github.com/isfk/gorm-cache/context"
	"gorm.io/gorm"
)

func GetDataFromCtx(ctx context.Context, stat *gorm.Statement) (key string, tags []string, valueMap map[string]any, err error) {
	keyValue := ctx.Value(cacheContext.GormCacheKeyCtx{})
	if keyValue != nil {
		if t, ok := keyValue.(string); ok {
			key = t
		}
	}

	tagsValue := ctx.Value(cacheContext.GormCacheTagsCtx{})
	if tagsValue != nil {
		if t, ok := tagsValue.([]string); ok {
			tags = t
		}
	}

	valsValue := ctx.Value(cacheContext.GormCacheValuesCtx{})

	var values []byte
	if valsValue != nil {
		if t, ok := valsValue.([]byte); ok {
			values = t
		}
	} else {
		values, err = json.Marshal(stat.Dest)
		if err != nil {
			return "", nil, nil, err
		}
	}

	if values != nil {
		err = json.Unmarshal(values, &valueMap)
		if err != nil {
			return "", nil, nil, err
		}
	}

	pf := strings.ToLower(stat.Schema.PrioritizedPrimaryField.Name)

	if strings.HasSuffix(key, ":") {
		key = getKeyFromMap(key, valueMap, pf)
	}
	for i, tag := range tags {
		if strings.HasSuffix(tag, ":") {
			tags[i] = getKeyFromMap(tag, valueMap, pf)
		}
	}

	if !NotEmpty(valueMap, pf) {
		valueMap = nil
	}

	return
}

func getKeyFromMap(key string, valueMap map[string]any, name string) string {
	switch valueMap[name].(type) {
	case float64:
		return fmt.Sprintf("%s%d", key, int64(valueMap[name].(float64)))
	case int64:
		return fmt.Sprintf("%s%d", key, valueMap[name].(int64))
	case string:
		return fmt.Sprintf("%s%s", key, valueMap[name].(string))
	}
	return ""
}

func NotEmpty(valueMap map[string]any, name string) bool {
	switch valueMap[name].(type) {
	case float64:
		id := int64(valueMap[name].(float64))
		if id > 0 {
			return true
		}
	case int64:
		id := valueMap[name].(int64)
		if id > 0 {
			return true
		}
	case string:
		id := valueMap[name].(string)
		if len(id) > 0 {
			return true
		}
	}
	return false
}
