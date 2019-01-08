build:
	docker build --rm -t redisgo .

run: build
	docker run -p 9090:9090/tcp --rm -it --name redisgo_1 redisgo

test:
	docker build --rm --target go_env -t redisgo-builder .
	docker run --rm -it redisgo-builder bash

check:
	docker build --rm --target go_env -t redisgo_go_env .
	docker run --rm -it redisgo_go_env ./run_check.sh
