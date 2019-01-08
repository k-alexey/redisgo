FROM golang:1.11 as go_env
WORKDIR /go/src/redisgo
COPY . .

FROM go_env as compiler
RUN CGO_ENABLED=0 go build -o client ./client.go
RUN CGO_ENABLED=0 go build -o server ./server.go

FROM alpine:latest
COPY --from=compiler /go/src/redisgo/client /usr/local/bin/
COPY --from=compiler /go/src/redisgo/server /usr/local/bin/
CMD [ "server"]
EXPOSE 9090
