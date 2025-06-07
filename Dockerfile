# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o db-backup .

FROM alpine:3.19

RUN apk add --no-cache \
    postgresql-client \
    mariadb-client \
    tzdata

RUN adduser -D -g '' appuser

WORKDIR /app

COPY --from=builder /app/db-backup .

COPY config.yaml .

RUN mkdir -p /backups && chown -R appuser:appuser /backups

USER appuser

ENV TZ=UTC

ENTRYPOINT ["/app/db-backup"]

CMD ["config.yaml"] 
