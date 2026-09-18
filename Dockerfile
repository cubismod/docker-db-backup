FROM golang:1.26-alpine@sha256:51a7c389a5ddaf82f527191a1e9bff9928655130a44e4975dd1d7e0acf59f1ae AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o db-backup .

FROM alpine:3.24@sha256:5b02b42e375f7426f8d65c3af331ca05d9878f9989230354504e0b9dfd431f60

RUN apk add --no-cache \
    postgresql-client \
    mariadb-client \
    redis \
    tzdata

WORKDIR /app

COPY --from=builder /app/db-backup .

COPY config.yaml .

ENV TZ=UTC

ENTRYPOINT ["/app/db-backup"]

CMD ["config.yaml"] 
