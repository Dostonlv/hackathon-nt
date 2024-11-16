.PHONY: run_db run test build

run_db:
	docker-compose up db -d

run:
	docker-compose up --build

test:
	go test ./...

build:
	go build -o main cmd/server/main.go

stop:
	docker-compose down

clean:
	docker-compose down -v