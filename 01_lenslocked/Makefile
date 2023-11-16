.PHONY: run
run:
	npm run start

.PHONY: status
status:
	cd migrations && goose postgres "host=localhost port=5432 user=goblina password=jinnythejimbo dbname=lenslocked sslmode=disable" status

.PHONY: up
up:
	cd migrations && goose postgres "host=localhost port=5432 user=goblina password=jinnythejimbo dbname=lenslocked sslmode=disable" up


.PHONY: down
down:
	cd migrations && goose postgres "host=localhost port=5432 user=goblina password=jinnythejimbo dbname=lenslocked sslmode=disable" down
