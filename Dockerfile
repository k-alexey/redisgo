FROM golang:1.11 as go_env
WORKDIR /go/src/redisgo
COPY . .

FROM go_env as compiler
RUN CGO_ENABLED=0 go install ./client
RUN CGO_ENABLED=0 go install ./server

FROM alpine:latest
COPY --from=compiler /go/bin/client /usr/local/bin/
COPY --from=compiler /go/bin/server /usr/local/bin/
CMD [ "server"]
EXPOSE 9090
