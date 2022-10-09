all: lint run-local

lint:
	golangci-lint run

docker-up:
	docker-compose up -d

docker-stop:
	docker-compose stop

run-local:
	docker-compose up -d
	go run ./src/main/main.go -config config-local.yml
