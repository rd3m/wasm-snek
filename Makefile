include .env
export REDIS_URL

build:
	docker build -t snek .

run:
	docker run -p 8080:8080 -e REDIS_URL=$$REDIS_URL snek

.PHONY: build run
