FROM golang:1.25-alpine@sha256:77dd832edf2752dafd030693bef196abb24dcba3a2bc3d7a6227a7a1dae73169 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o db-backup .

FROM alpine:3.22@sha256:4bcff63911fcb4448bd4fdacec207030997caf25e9bea4045fa6c8c44de311d1

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
