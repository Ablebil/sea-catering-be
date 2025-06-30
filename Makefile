run:
	@docker compose up
build:
	@docker compose down
	@docker compose up --build
stop:
	@docker compose down
logs:
	@docker compose  logs -f
restart:
	@docker compose down
	@docker compose up
down-remove-volumes:
	@docker compose down -v