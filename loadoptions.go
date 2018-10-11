package pantry

import (
	"context"
)

type Watcher func(err error)
type LoadOptions struct {
	Watcher Watcher
	Context context.Context
	Reload  bool
}

func Watch(w Watcher) func(*LoadOptions) {
	return func(lo *LoadOptions) {
		lo.Watcher = w
		lo.Context = context.Background()
	}
}

func WatchContext(w Watcher, ctx context.Context) func(*LoadOptions) {
	return func(lo *LoadOptions) {
		lo.Watcher = w
		lo.Context = ctx
	}
}

func AutoReload() func(*LoadOptions) {
	return func(lo *LoadOptions) {
		lo.Reload = true
	}
}
