FROM golang:1.27-alpine AS builder

WORKDIR /src

# goose is used at startup to apply migrations.
RUN go install github.com/pressly/goose/v3/cmd/goose@v3.24.1

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o /out/eventer .

FROM alpine:3.20

RUN apk add --no-cache ca-certificates \
    && addgroup -S app \
    && adduser -S -G app app

WORKDIR /app

COPY --from=builder /out/eventer /app/eventer
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY migrations /app/migrations
COPY docker-entrypoint.sh /app/docker-entrypoint.sh
RUN chmod +x /app/docker-entrypoint.sh

USER app

EXPOSE 8080

HEALTHCHECK --interval=10s --timeout=3s --retries=5 \
    CMD wget -qO- http://localhost:8080/healthz || exit 1

ENTRYPOINT ["/app/docker-entrypoint.sh"]
