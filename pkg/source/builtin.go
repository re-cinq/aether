package source

import (
	"context"

	"github.com/re-cinq/aether/pkg/config"
	amazon "github.com/re-cinq/aether/pkg/providers/aws"
	"github.com/re-cinq/aether/pkg/providers/gcp"
	v1 "github.com/re-cinq/aether/pkg/types/v1"
)

func BuiltInSources(ctx context.Context) []v1.Source {
	var sources []v1.Source

	if cfg, exists := config.AppConfig().Providers[v1.GCP]; exists {
		sources = append(sources, gcp.Sources(ctx, &cfg)...)
	}

	if cfg, exists := config.AppConfig().Providers[v1.AWS]; exists {
		sources = append(sources, amazon.Sources(ctx, &cfg)...)
	}

	return sources
}
