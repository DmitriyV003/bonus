start:
	docker-compose up --build
migrate:
	tern migrate --migrations ./cmd/gophermart/migrations