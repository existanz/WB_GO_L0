package usecases

import (
	"context"
	"fmt"

	"WB_GO_L0/internal/entities"
	"WB_GO_L0/internal/kafka"

	"github.com/brianvoe/gofakeit/v7"
)

func GenerateAndSendMessages(ctx context.Context, n int) error {
	if n <= 0 || n > 100 {
		n = 100
	} // Limit the number of messages to 100

	for i := 0; i < n; i++ {
		ord, err := generateStruct()
		if err != nil {
			return err
		}

		err = kafka.SendMessage(ctx, &ord)
		if err != nil {
			return err
		}

		fmt.Printf("\n Sent message: %s", ord.String())
	}

	return nil
}

func generateStruct() (entities.Order, error) {
	ord := entities.Order{}
	err := gofakeit.Struct(&ord)
	return ord, err
}
