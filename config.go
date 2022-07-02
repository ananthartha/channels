package channels

import (
	"context"

	"go.opentelemetry.io/otel/metric"
)

type ConfigFn[I, O any] func(*Config[I, O]) error

type Config[I, O any] struct {
	Context         context.Context
	Channel         UnbufferedChannel[I]
	InitTaskManager func(*TaskManagerContext[I, O]) TaskManager[I]
}

func defaultConfig[I, O any]() *Config[I, O] {
	instance := &Config[I, O]{
		Context: context.Background(),
		//TODO: Later Changes this to blocking Queue
		Channel: NewQueueChannel[I](NewDefaultQueue[I]()),
	}

	WithRestrictedTaskManager[I, O](10)(instance)
	return instance
}

func buildConfig[I, O any](configs ...ConfigFn[I, O]) *Config[I, O] {
	config := defaultConfig[I, O]()

	for _, configFn := range configs {
		if configFn != nil {
			configFn(config)
		}
	}

	return config
}

func WithContext[I, O any](ctx context.Context) ConfigFn[I, O] {
	return func(c *Config[I, O]) error {
		c.Context = ctx
		return nil
	}
}

func WithChannel[I, O any](channel Channel[I]) ConfigFn[I, O] {
	return func(c *Config[I, O]) error {
		c.Channel = channel
		return nil
	}
}

func WithChannelAndMeterics[I, O any](channel Channel[I], prefix string, meter metric.Meter) ConfigFn[I, O] {
	return func(c *Config[I, O]) error {
		c.Channel = channel
		RegisterBuffer(prefix, meter, channel)
		return nil
	}
}

// Will be blocking channels
func WithUnbufferedChannel[I, O any](channel UnbufferedChannel[I]) ConfigFn[I, O] {
	return func(c *Config[I, O]) error {
		c.Channel = channel
		return nil
	}
}
