# Build stage
FROM golang:1.22 AS builder

WORKDIR /app

# Copy source code
COPY . .

# Download dependencies and build
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o echo-api-clean-architecture

# Final stage
FROM alpine:latest 

WORKDIR /app

# Install necessary packages
RUN apk add --no-cache ca-certificates tzdata

# Copy binary and config
COPY --from=builder /app/echo-api-clean-architecture .
COPY .env.prod .

EXPOSE 3000

# Run the application
ENTRYPOINT [ "./echo-api-clean-architecture", "./.env.prod" ]