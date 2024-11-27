package kafka

import (
	"context"
	"fmt"

	"WB_GO_L0/internal/entity"

	"github.com/segmentio/kafka-go"
)

func SendMessage(ctx context.Context, msg *entity.Order) error {
	w := &kafka.Writer{
		Addr:  kafka.TCP(brokers...),
		Topic: topic,
	}

	defer w.Close()

	messageBytes := []byte(msg.String())

	err := w.WriteMessages(ctx, kafka.Message{
		Key:   []byte(msg.OrderUID),
		Value: messageBytes,
	})
	if err != nil {
		return fmt.Errorf("error writing message to kafka: %w", err)
	}

	return nil
}
