package channels

import (
	"context"

	"go.opentelemetry.io/otel/metric"
)

type ConfigFn[I, O any] func(*config[I, O]) error

type config[I, O any] struct {
	context         context.Context
	channel         UnbufferedChannel[I]
	initTaskManager func(*TaskManagerContext[I, O]) TaskManager[I]
}

func defaultConfig[I, O any]() *config[I, O] {
	instance := &config[I, O]{
		context: context.Background(),
		//TODO: Later Changes this to blocking Queue
	}

	WithRestrictedTaskManager[I, O](10)(instance)
	return instance
}

func buildConfig[I, O any](configs ...ConfigFn[I, O]) *config[I, O] {
	config := defaultConfig[I, O]()

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

func WithContext[I, O any](ctx context.Context) ConfigFn[I, O] {
	return func(c *config[I, O]) error {
		c.context = ctx
		return nil
	}
}

func WithChannel[I, O any](channel Channel[I], prefix string, meter metric.Meter) ConfigFn[I, O] {
	return func(c *config[I, O]) error {
		c.channel = channel
		if meter != nil {
			RegisterBuffer(prefix, meter, channel)
		}
		return nil
	}
}

// Will be blocking channels
func WithUnbufferedChannel[I, O any](channel UnbufferedChannel[I]) ConfigFn[I, O] {
	return func(c *config[I, O]) error {
		c.channel = channel
		return nil
	}
}
