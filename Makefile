build:
	docker build --rm -t redisgo .

run: run_server

run_server: build
	docker run -p 9090:9090/tcp --rm -it --name redisgo_1 redisgo

run_client:
	docker exec -it redisgo_1 client

test:
	docker build --rm --target go_env -t redisgo-go_env .
	docker run -it --name redisgo_testing redisgo-go_env ./run_test.sh
	docker cp redisgo_testing:/go/src/redisgo/coverage.out ./coverage.out
	docker container rm redisgo_testing

check:
	docker build --rm --target go_env -t redisgo_go_env .
	docker run --rm -it redisgo_go_env ./run_check.sh
