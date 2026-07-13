package events

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

// EnsureTopic creates a Kafka topic if it does not already exist.
// Safe to call on every service start (idempotent).
func EnsureTopic(brokersCSV, topic string, partitions int) error {
	brokers := splitBrokers(brokersCSV)
	if len(brokers) == 0 {
		return fmt.Errorf("kafka brokers required")
	}
	topic = strings.TrimSpace(topic)
	if topic == "" {
		return fmt.Errorf("kafka topic required")
	}
	if partitions <= 0 {
		partitions = 1
	}

	conn, err := kafka.Dial("tcp", brokers[0])
	if err != nil {
		return fmt.Errorf("dial kafka broker %s: %w", brokers[0], err)
	}
	defer func() { _ = conn.Close() }()

	_ = conn.SetDeadline(time.Now().Add(10 * time.Second))
	controller, err := conn.Controller()
	if err != nil {
		return fmt.Errorf("kafka controller: %w", err)
	}

	controllerAddr := net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port))
	controllerConn, err := kafka.Dial("tcp", controllerAddr)
	if err != nil {
		return fmt.Errorf("dial kafka controller %s: %w", controllerAddr, err)
	}
	defer func() { _ = controllerConn.Close() }()

	_ = controllerConn.SetDeadline(time.Now().Add(10 * time.Second))
	err = controllerConn.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     partitions,
		ReplicationFactor: 1,
	})
	if err != nil && !isTopicExists(err) {
		return fmt.Errorf("create topic %s: %w", topic, err)
	}
	return nil
}

func splitBrokers(brokersCSV string) []string {
	parts := strings.Split(brokersCSV, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func isTopicExists(err error) bool {
	if err == nil {
		return false
	}
	// kafka-go returns TopicAlreadyExists among other wrapped forms
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "already exists") || strings.Contains(msg, "topic with this name already exists")
}
