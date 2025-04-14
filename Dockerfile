# Stage 1: Build Stage
FROM golang:alpine AS builder

# Set build environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Set working directory for building
WORKDIR /app

# Copy dependency definitions to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application binary
RUN go build -o main .

# Stage 2: Production Stage
FROM arillso/ansible:2.18.4

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main /usr/local/bin/main

# Set the production user to "ansible"
USER ansible

# Set the default entrypoint to run the application
ENTRYPOINT ["/usr/local/bin/main"]

# Healthcheck to verify Ansible functionality
HEALTHCHECK --interval=30s --timeout=10s CMD ansible --version || exit 1
