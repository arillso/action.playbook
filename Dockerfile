# -------------------------
# Stage 1: Build Stage
# -------------------------
FROM golang:alpine@sha256:aee43c3ccbf24fdffb7295693b6e33b21e01baec1b2a55acc351fde345e9ec34 AS builder
# Use the official Golang Alpine image and assign this build stage the name "builder".

# Set build environment variables for module support and cross-compilation.
ENV GO111MODULE=on \
	CGO_ENABLED=0 \
	GOOS=linux \
	GOARCH=amd64

# Set the working directory inside the container where the build will be executed.
WORKDIR /app

# Copy dependency definitions (go.mod and go.sum) to leverage Docker's caching mechanism.
COPY go.mod go.sum ./

# Download all Go module dependencies.
RUN go mod download

# Copy the entire source code into the container.
COPY . .

# Build the application binary and output it as "main".
RUN go build -o main .

# -------------------------
# Stage 2: Production Stage (Alpine Linux)
# -------------------------
FROM arillso/ansible:2.20.0@sha256:727a1be9ec328c3909a190312dd9764e6981881e2acaf076606758f21820e9da
# Use an Ansible-based Alpine Linux image as the base for the production stage.

# Switch to root user to execute system-level modifications.
USER root

# Update package repositories and install the 'shadow' package for access to usermod and groupmod commands.
RUN apk update && \
	apk add --no-cache \
	shadow=4.17.3-r0

# Modify the UID and GID of the 'ansible' user from 1000 to 1001 and update file ownership:
# - Change the UID of user 'ansible' to 1001.
# - Change the GID of group 'ansible' to 1001.
# - Recursively change ownership of files from the old UID (1000) to the new UID (1001).
RUN usermod -u 1001 ansible && \
	groupmod -g 1001 ansible && \
	find / -xdev -user 1000 -exec chown -h 1001 {} \;

# Copy the compiled binary from the builder stage into the final image's binary directory.
COPY --from=builder /app/main /usr/local/bin/main

# Switch to the non-privileged 'ansible' user for runtime.
USER ansible

# Set a working directory for the 'ansible' user.
WORKDIR /home/ansible

# Set the default entrypoint to execute the application binary.
ENTRYPOINT ["/usr/local/bin/main"]

# Add a healthcheck to ensure Ansible functionality:
# It runs 'ansible --version' every 30 seconds with a timeout of 10 seconds. If the command fails, the container exits with status 1.
HEALTHCHECK --interval=30s --timeout=10s CMD ansible --version || exit 1
