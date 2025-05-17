package bucket

import "context"

type Limiter interface {
	Allow(ctx context.Context, key string) (bool, error)
	Reset(ctx context.Context, key string) error
	Close(ctx context.Context) error
}
