FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/scheduler .

FROM alpine:latest

ENV TODO_PORT=7540
ENV TODO_DBFILE=/data/scheduler.db

WORKDIR /app

RUN mkdir -p /data

COPY --from=builder /app/scheduler /app/scheduler
COPY --from=builder /app/web /app/web

EXPOSE 7540

ENTRYPOINT ["/app/scheduler"]