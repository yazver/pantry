package pantry

import (
	"context"
)

type Watcher func(err error)
type LoadOptions struct {
	Watcher Watcher
	Context context.Context
	Reload  bool
	Format  string
}

func Watch(w Watcher) LoadOptions {
	return LoadOptions{Watcher: w, Context: context.Background(), Reload: false}
}

func WatchContext(w Watcher, ctx context.Context) LoadOptions {
	return LoadOptions{Watcher: w, Context: ctx, Reload: false}
}

func Reload(w Watcher) LoadOptions {
	return LoadOptions{Watcher: w, Context: context.Background(), Reload: true}
}

func ReloadContext(w Watcher, ctx context.Context) LoadOptions {
	return LoadOptions{Watcher: w, Context: ctx, Reload: true}
}

func Format(s string) LoadOptions {
	return LoadOptions{Format: s}
}
