FROM golang:1.23.4-alpine AS builder

WORKDIR /app
COPY go.* ./
COPY . .

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Download dependencies
RUN go mod download

# Build the application
RUN CGO_ENABLED=1 go build -o main .

# Final stage
FROM alpine:3.19

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates

# Copy the binary from builder
COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
