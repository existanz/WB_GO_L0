package kafka

import (
	"context"
	"log"
	"os"
	"strings"

	_ "github.com/joho/godotenv/autoload" // Load the .env file
	"github.com/segmentio/kafka-go"
)

var (
	brokers = os.Getenv("BROKERS")
	topic   = os.Getenv("TOPIC")
)

func Consume(ctx context.Context) error {
	brokersArr := strings.Split(brokers, ",")
	return consume(ctx, brokersArr, topic)
}

func consume(ctx context.Context, brokers []string, topic string) error {
	brokersArr := brokers
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokersArr,
		Topic:   topic,
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
			log.Printf("message at topic:%v partition:%v offset:%v	key:%s value:%s", msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value) // Process the message here
		}
	}
}
