FROM golang:1.18.1

COPY ./ ./
RUN go build -o main .
CMD ["./main"]