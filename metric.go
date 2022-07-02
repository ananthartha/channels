package channels

import (
	context "context"
	"fmt"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument"
)

func RegisterBuffer(prefix string, meter metric.Meter, buffer Buffer) {
	var instrumentProvider = meter.AsyncInt64()

	bufferLength, err := instrumentProvider.Gauge(fmt.Sprintf("%s_length", prefix))
	if err != nil {
		// logger.Panicf("failed to initialize tickEnqueued: %v", err)
		return
	}

	bufferCapacity, err := instrumentProvider.Gauge(fmt.Sprintf("%s_capacity", prefix))
	if err != nil {
		// logger.Panicf("failed to initialize tickEnqueued: %v", err)
		return
	}

	_ = meter.RegisterCallback([]instrument.Asynchronous{bufferLength, bufferCapacity}, func(ctx context.Context) {
		bufferLength.Observe(ctx, int64(buffer.Len()))
		bufferCapacity.Observe(ctx, int64(buffer.Cap()))
	})
}

func RegisterTaskManager[I any](prefix string, meter metric.Meter, taskManager TaskManager[I]) {
	var instrumentProvider = meter.AsyncInt64()

	taskManagerSize, err := instrumentProvider.Gauge(fmt.Sprintf("%s_size", prefix))
	if err != nil {
		// logger.Panicf("failed to initialize tickEnqueued: %v", err)
		return
	}

	taskManagerActive, err := instrumentProvider.Gauge(fmt.Sprintf("%s_active", prefix))
	if err != nil {
		// logger.Panicf("failed to initialize tickEnqueued: %v", err)
		return
	}

	taskManagerThreshold, err := instrumentProvider.Gauge(fmt.Sprintf("%s_threshold", prefix))
	if err != nil {
		// logger.Panicf("failed to initialize tickEnqueued: %v", err)
		return
	}

	taskManagerSkipped, err := instrumentProvider.Gauge(fmt.Sprintf("%s_skipped", prefix))
	if err != nil {
		// logger.Panicf("failed to initialize tickEnqueued: %v", err)
		return
	}

	_ = meter.RegisterCallback([]instrument.Asynchronous{taskManagerSize, taskManagerActive, taskManagerThreshold, taskManagerSkipped}, func(ctx context.Context) {
		taskManagerSize.Observe(ctx, int64(taskManager.Size()))
		taskManagerActive.Observe(ctx, int64(taskManager.Active()))
		taskManagerThreshold.Observe(ctx, int64(taskManager.Threshold()))
		taskManagerSkipped.Observe(ctx, int64(taskManager.Skipped()))
	})
}
