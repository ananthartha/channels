package channels

import (
	"context"

	"go.opentelemetry.io/otel/metric"
)

type ConfigFn[I any] func(*config[I]) error

type config[I any] struct {
	context         context.Context
	channel         UnbufferedChannel[I]
	initTaskManager func(*TaskManagerContext[I]) TaskManager[I]
}

func defaultConfig[I any]() *config[I] {
	instance := &config[I]{
		context: context.Background(),
		//TODO: Later Changes this to blocking Queue
	}

	WithRestrictedTaskManager[I](10)(instance)
	return instance
}

func buildConfig[I any](configs ...ConfigFn[I]) *config[I] {
	config := defaultConfig[I]()

	for _, configFn := range configs {
		if configFn != nil {
			configFn(config)
		}
	}

	if config.channel == nil {
		config.channel = NewQueueChannel[I](NewDefaultQueue[I]())
	}

	return config
}

func WithContext[I, O any](ctx context.Context) ConfigFn[I] {
	return func(c *config[I]) error {
		c.context = ctx
		return nil
	}
}

func WithChannel[I, O any](channel Channel[I], prefix string, meter metric.Meter) ConfigFn[I] {
	return func(c *config[I]) error {
		c.channel = channel
		if meter != nil {
			RegisterBuffer(prefix, meter, channel)
		}
		return nil
	}
}

// Will be blocking channels
func WithUnbufferedChannel[I, O any](channel UnbufferedChannel[I]) ConfigFn[I] {
	return func(c *config[I]) error {
		c.channel = channel
		return nil
	}
}
