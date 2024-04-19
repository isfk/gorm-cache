package context

import "context"

type (
	GormCacheKeyCtx    struct{}
	GormCacheTagsCtx   struct{}
	GormCacheValuesCtx struct{}
)

type CacheContext struct {
	Ctx context.Context
}

func New(ctx context.Context) *CacheContext {
	return &CacheContext{
		Ctx: ctx,
	}
}

func (c *CacheContext) WithKey(key string) *CacheContext {
	c.Ctx = context.WithValue(c.Ctx, GormCacheKeyCtx{}, key)
	return c
}

func (c *CacheContext) WithTags(tags ...string) *CacheContext {
	c.Ctx = context.WithValue(c.Ctx, GormCacheTagsCtx{}, tags)
	return c
}

func (c *CacheContext) CC() context.Context {
	return c.Ctx
}
