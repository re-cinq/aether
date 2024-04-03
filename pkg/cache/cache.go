package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	bc "github.com/allegro/bigcache/v3"
	"github.com/eko/gocache/lib/v4/marshaler"
	store "github.com/eko/gocache/store/bigcache/v4"
	"github.com/re-cinq/aether/pkg/config"
)

func New(ctx context.Context) (*marshaler.Marshaler, error) {
	cfg := config.AppConfig()
	// Check the config is set, otherwise tests will panic
	if cfg == nil {
		return nil, errors.New("config not set")
	}

	switch cfg.Cache.Store {
	case store.BigcacheType:
		return bigcache(ctx, cfg.Cache.Expiry)
	default:
		return nil, fmt.Errorf("error cache not yet supported: %s", cfg.Cache.Store)
	}
}

func bigcache(ctx context.Context, expiry time.Duration) (*marshaler.Marshaler, error) {
	cli, err := bc.New(ctx, bc.DefaultConfig(expiry))
	if err != nil {
		return nil, err
	}

	return marshaler.New(store.NewBigcache(cli)), nil
}
