.PHONY: run
run:
	npm run start

.PHONY: status
status:
	goose --dir ./migrations postgres "host=localhost port=5432 user=goblina password=jinnythejimbo dbname=lenslocked sslmode=disable" status

.PHONY: up
up:
	goose --dir ./migrations postgres "host=localhost port=5432 user=goblina password=jinnythejimbo dbname=lenslocked sslmode=disable" up

.PHONY: down
down:
	goose --dir ./migrations postgres "host=localhost port=5432 user=goblina password=jinnythejimbo dbname=lenslocked sslmode=disable" down

.PHONY: reset
reset:
	goose --dir ./migrations postgres "host=localhost port=5432 user=goblina password=jinnythejimbo dbname=lenslocked sslmode=disable" reset

.PHONY: fix
fix:
	goose --dir ./migrations fix

.PHONY: updocker
updocker:
	sudo docker compose up -d

.PHONY: downducker
downducker:
	sudo docker compose down
