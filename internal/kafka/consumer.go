package kafka

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"os"
	"strings"

	"WB_GO_L0/internal/database"

	_ "github.com/joho/godotenv/autoload" // Load the .env file
	"github.com/segmentio/kafka-go"
)

var (
	brokers = strings.Split(os.Getenv("BROKERS"), ",")
	topic   = os.Getenv("TOPIC")
)

func Consume(ctx context.Context, db database.Service) error {
	if len(brokers) == 0 || topic == "" {
		return errors.New("BROKERS and TOPIC environment variables must be set")
	}
	return consume(ctx, db, brokers, topic)
}

func consume(ctx context.Context, db database.Service, brokers []string, topic string) error {
	brokersArr := brokers
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokersArr,
		Topic:   topic,
		GroupID: "group1",
	})
	defer reader.Close()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			msg, err := reader.ReadMessage(ctx)
			if err != nil {
				return err
			}
			slog.Info("message at ", "topic", msg.Topic, "partition", msg.Partition, "offset", msg.Offset, "key", msg.Key, "value", msg.Value)

			err = db.SaveOrderPlain(string(msg.Value))
			if err != nil {
				log.Fatal(err)
			}
			if err := reader.CommitMessages(ctx, msg); err != nil {
				log.Fatal("Failed to commit message:", err)
			}
		}
	}
}
