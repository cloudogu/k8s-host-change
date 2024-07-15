package alias

import (
	"context"
	"github.com/cloudogu/k8s-registry-lib/config"
)

type globalConfigGetter interface {
	Get(ctx context.Context) (config.GlobalConfig, error)
}
