package kafka

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"time"

	workflowapp "github.com/ai-finops/ai-finops/services/workflow-service/internal/application/workflow"
	"github.com/ai-finops/ai-finops/packages/events"
	ingestionv1 "github.com/ai-finops/ai-finops/packages/proto/gen/go/ingestion/v1"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

const maxProcessAttempts = 3

// Consumer reads usage event messages and dispatches them to the workflow service.
type Consumer struct {
	reader *kafka.Reader
	svc    *workflowapp.Service
	log    *slog.Logger
}

// NewConsumer creates a Kafka consumer for usage.events.raw.
func NewConsumer(brokers, topic, groupID string, svc *workflowapp.Service, log *slog.Logger) *Consumer {
	brokerList := strings.Split(brokers, ",")
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:               brokerList,
			Topic:                 topic,
			GroupID:               groupID,
			MinBytes:              1,
			MaxBytes:              10e6,
			CommitInterval:        time.Second,
			StartOffset:           kafka.FirstOffset,
			WatchPartitionChanges: true,
			MaxWait:               time.Second,
		}),
		svc: svc,
		log: log,
	}
}

// Run consumes messages until the context is cancelled.
func (c *Consumer) Run(ctx context.Context) error {
	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil || err == io.EOF {
				return nil
			}
			c.log.Error("fetch kafka message failed", slog.String("error", err.Error()))
			continue
		}

		if err := c.handleMessage(ctx, msg.Value); err != nil {
			c.log.Error("process kafka message failed",
				slog.String("topic", msg.Topic),
				slog.Int("partition", msg.Partition),
				slog.Int64("offset", msg.Offset),
				slog.String("error", err.Error()),
			)
			continue
		}

		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			c.log.Error("commit kafka offset failed", slog.String("error", err.Error()))
		}
	}
}

func (c *Consumer) handleMessage(ctx context.Context, value []byte) error {
	var created ingestionv1.UsageEventCreated
	if err := proto.Unmarshal(value, &created); err != nil {
		return fmt.Errorf("unmarshal usage event created: %w", err)
	}

	var lastErr error
	for attempt := 1; attempt <= maxProcessAttempts; attempt++ {
		_, err := c.svc.ProcessUsageEventCreated(ctx, &created)
		if err == nil {
			return nil
		}
		lastErr = err
		c.log.Warn("retrying usage event processing",
			slog.String("event_id", created.GetEventId()),
			slog.Int("attempt", attempt),
			slog.String("error", err.Error()),
		)
		time.Sleep(time.Duration(attempt) * 200 * time.Millisecond)
	}
	return fmt.Errorf("process usage event after %d attempts: %w", maxProcessAttempts, lastErr)
}

// Close closes the Kafka reader.
func (c *Consumer) Close() error {
	if c.reader == nil {
		return nil
	}
	c.log.Info("kafka consumer stopping", slog.String("topic", events.TopicUsageEventsRaw))
	return c.reader.Close()
}