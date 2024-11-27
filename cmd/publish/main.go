package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"WB_GO_L0/internal/usecase"
)

func main() {
	var n int
	flag.IntVar(&n, "n", 1, "number of messages to generate and send")
	flag.Parse()

	fmt.Printf("Generating and sending %d messages...\n", n)
	err := usecase.GenerateAndSendMessages(context.Background(), n)
	if err != nil {
		log.Fatal(err)
	}
}
