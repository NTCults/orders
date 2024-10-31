run:
	go run ./cmd/main.go

compose:
	docker-compose -f ./build/docker-compose.yml up --build