start:
	docker-compose up --build
migrate:
	tern migrate --migrations ./migrations