# Builder stage
FROM golang:alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o zama .

# Runtime stage
FROM alpine:latest AS runtime

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/zama .

ADD test test

# Copy any other files you need at runtime (e.g., configuration files, certificates)
# COPY --from=builder /app/configs ./configs

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./zama"]
