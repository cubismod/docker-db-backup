FROM golang:1.25-alpine@sha256:f6751d823c26342f9506c03797d2527668d095b0a15f1862cddb4d927a7a4ced AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o db-backup .

FROM alpine:3.23@sha256:25109184c71bdad752c8312a8623239686a9a2071e8825f20acb8f2198c3f659

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
