FROM golang:alpine as builder

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY . .

RUN go mod download

# Build the application
RUN go build -o main

FROM ubuntu:latest as production

RUN apt install python3 \
    pip3 install --upgrade pip3 \ 
    pip3 install boto3 botocore \
    pip3 install ansible ansible-core

# Copy binary from build to main folder
COPY --from=builder /build/main /usr/local/bin

# Run as root
USER root

# Command to run when starting the container
CMD ["main"]
