compose-up:
	docker compose up --build -d && docker compose logs -f

lint:
	golangci-lint run

run:
	go run main.go
