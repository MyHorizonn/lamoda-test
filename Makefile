.PHONY: up test
up:
	docker-compose build && docker-compose up
test:
	go test ./...