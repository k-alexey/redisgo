build:
	docker build --rm -t redisgo .

run: build
	docker run -p 9090:9090/tcp --rm -it redisgo

test:
	docker build --rm --target builder -t redisgo-builder .
	docker run --rm -it redisgo-builder bash

check:
	docker build --rm --target builder -t redisgo-builder .
	docker run --rm -it redisgo-builder ./run_check.sh
