package usecase

import (
	"context"
	"fmt"

	"WB_GO_L0/internal/entity"
	"WB_GO_L0/internal/kafka"

	"github.com/brianvoe/gofakeit/v7"
)

func GenerateAndSendMessages(ctx context.Context, n int) error {
	if n <= 0 || n > 100 {
		n = 100
	} // Limit the number of messages to 100

	for i := 0; i < n; i++ {
		ord, err := GenerateStruct()
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

func GenerateStruct() (entity.Order, error) {
	ord := entity.Order{}
	err := gofakeit.Struct(&ord)
	return ord, err
}
