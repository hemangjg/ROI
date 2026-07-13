package relay

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/ai-finops/ai-finops/packages/events"
	svcconfig "github.com/ai-finops/ai-finops/services/outbox-relay/internal/config"
	sqlcgen "github.com/ai-finops/ai-finops/db/atlas/sql/gen"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

// Relay polls unpublished outbox rows and publishes them to Kafka.
type Relay struct {
	pool     *pgxpool.Pool
	queries  *sqlcgen.Queries
	writer   *kafka.Writer
	interval time.Duration
	batch    int32
	log      *slog.Logger
}

// New creates an outbox relay.
func New(pool *pgxpool.Pool, queries *sqlcgen.Queries, cfg svcconfig.Config, log *slog.Logger) (*Relay, error) {
	interval, err := time.ParseDuration(cfg.PollInterval)
	if err != nil {
		return nil, fmt.Errorf("parse poll_interval: %w", err)
	}

	brokers := strings.Split(cfg.KafkaBrokers, ",")
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(brokers...),
		Topic:                  cfg.KafkaTopic,
		Balancer:               &kafka.Hash{},
		AllowAutoTopicCreation: true,
	}

	return &Relay{
		pool:     pool,
		queries:  queries,
		writer:   writer,
		interval: interval,
		batch:    int32(cfg.BatchSize),
		log:      log,
	}, nil
}

// Run polls until the context is cancelled.
func (r *Relay) Run(ctx context.Context) error {
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()

	for {
		if err := r.pollOnce(ctx); err != nil {
			r.log.Error("outbox poll failed", slog.String("error", err.Error()))
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}

func (r *Relay) pollOnce(ctx context.Context) error {
	rows, err := r.queries.ListUnpublishedOutboxEvents(ctx, r.batch)
	if err != nil {
		return fmt.Errorf("list unpublished outbox events: %w", err)
	}
	if len(rows) == 0 {
		return nil
	}

	for _, row := range rows {
		if err := r.publishRow(ctx, row); err != nil {
			r.log.Error("publish outbox event failed",
				slog.String("error", err.Error()),
			)
			continue
		}
	}
	return nil
}

func (r *Relay) publishRow(ctx context.Context, row sqlcgen.OutboxEvent) error {
	message, err := payloadToProto(row.EventType, row.Payload)
	if err != nil {
		return err
	}

	body, err := proto.Marshal(message)
	if err != nil {
		return fmt.Errorf("marshal usage event created: %w", err)
	}

	if err := r.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(message.GetOrgId()),
		Value: body,
	}); err != nil {
		return fmt.Errorf("write kafka message: %w", err)
	}

	if err := r.queries.MarkOutboxEventPublished(ctx, row.ID); err != nil {
		return fmt.Errorf("mark outbox event published: %w", err)
	}

	r.log.Info("published outbox event",
		slog.String("event_id", message.GetEventId()),
		slog.String("org_id", message.GetOrgId()),
		slog.String("topic", events.TopicUsageEventsRaw),
	)
	return nil
}

// Close releases relay resources.
func (r *Relay) Close() error {
	if r.writer == nil {
		return nil
	}
	return r.writer.Close()
}