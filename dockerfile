FROM golang:1.18-alpine as compiler

WORKDIR /app

COPY . .

RUN go build -o /app/golang-project

FROM alpine:latest

WORKDIR /small

COPY --from=compiler /app/golang-project ./binary

RUN apk --no-cache add ca-certificates && \
    rm -rf /var/cache/apk/* /tmp/*

ENTRYPOINT [ "./binary" ]
