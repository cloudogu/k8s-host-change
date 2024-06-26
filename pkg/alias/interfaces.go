package alias

import "context"

type globalConfigValueGetter interface {
	Get(ctx context.Context, key string) (string, error)
	GetAll(ctx context.Context) (map[string]string, error)
}
