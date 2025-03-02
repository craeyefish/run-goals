# Build Stage
FROM golang:1.23-alpine AS builder
WORKDIR /app

# Install gcc and musl-dev for CGO support
RUN apk add --no-cache gcc musl-dev

# Copy go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go binary with CGO enabled
RUN CGO_ENABLED=1 go build -o backend .

# Final Stage: Minimal container with the binary
FROM alpine:latest
WORKDIR /app
# Copy the binary from the builder
COPY --from=builder /app/backend .
EXPOSE 8080
ENTRYPOINT ["./backend"]
