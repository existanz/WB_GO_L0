package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"log/slog"
	"os"
	"strings"

	"WB_GO_L0/internal/database"
	"WB_GO_L0/internal/entities"

	"github.com/go-playground/validator/v10"
	_ "github.com/joho/godotenv/autoload" // Load the .env file
	"github.com/segmentio/kafka-go"
)

var (
	brokers = strings.Split(os.Getenv("BROKERS"), ",")
	topic   = os.Getenv("TOPIC")
)

type KafkaOrderConsumer struct {
	r         *kafka.Reader
	db        database.Service
	validator *validator.Validate
}

func NewKafkaOrderConsumer(db database.Service) *KafkaOrderConsumer {
	if len(brokers) == 0 || topic == "" {
		log.Fatal(errors.New("BROKERS and TOPIC environment variables must be set"))
	}
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: "group1",
	})
	return &KafkaOrderConsumer{
		r:         r,
		db:        db,
		validator: validator.New(),
	}
}

func (k *KafkaOrderConsumer) Consume(ctx context.Context) error {
	defer k.r.Close()
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			msg, err := k.r.ReadMessage(ctx)
			if err != nil {
				return err
			}
			slog.Info("message at ", "topic", msg.Topic, "partition", msg.Partition, "offset", msg.Offset, "key", string(msg.Key))

			k.processMessageValue(msg.Value)

			if err := k.r.CommitMessages(ctx, msg); err != nil {
				slog.Error(err.Error())
				return err
			}
		}
	}
}

func (k *KafkaOrderConsumer) processMessageValue(value []byte) {
	var order entities.Order
	err := json.Unmarshal(value, &order)
	if err != nil {
		slog.Debug("Order is not valid json", "order", string(value))
		slog.Error(err.Error())
		return
	}
	if err := k.validator.Struct(order); err != nil {
		slog.Debug("Invalid order", "order", order)
		slog.Error(err.Error())
		return
	}
	slog.Info("Validation passed", "order uid", order.OrderUID)
	slog.Info("Saving order to database...")
	err = k.db.SaveOrderPlain(string(value))
	if err != nil {
		slog.Error(err.Error())
	}
}
