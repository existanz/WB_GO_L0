all: build test

build:
	@echo "Building..."
	@go build -o main cmd/api/main.go

run:
	@go run cmd/api/main.go

send:
	@go run cmd/publish/main.go -n=3

db:
	docker compose up --build

test:
	@echo "Testing..."
	@go test ./... -v

clean:
	@echo "Cleaning..."
	@rm -f main