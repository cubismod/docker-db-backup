# Build stage
FROM golang:1.21-alpine@sha256:2414035b086e3c42b99654c8b26e6f5b1b1598080d65fd03c7f499552ff4dc94 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o db-backup .


FROM alpine:3.22@sha256:8a1f59ffb675680d47db6337b49d22281a139e9d709335b492be023728e11715

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
