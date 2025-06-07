# Build stage
FROM golang:1.24-alpine@sha256:68932fa6d4d4059845c8f40ad7e654e626f3ebd3706eef7846f319293ab5cb7a AS builder

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
