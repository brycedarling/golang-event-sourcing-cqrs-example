FROM golang:1.15.2 AS builder

RUN mkdir /app
ADD . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/server/main.go

FROM alpine:latest
COPY --from=builder /app .
CMD ["./main"]
