# Start from the latest golang base image
FROM golang:1.22.5-alpine3.20 AS builder

ENV PATH /usr/local/go/bin:$PATH
ENV GOLANG_VERSION 1.22.5


# Add Maintainer Info
LABEL maintainer="cgil"
LABEL org.opencontainers.image.title="go-cloud-k8s-jwt-login"
LABEL org.opencontainers.image.description="This is a go-cloud-k8s-jwt-login container image, Allows to get a jwt token from a userid received by header (like f5)"
LABEL org.opencontainers.image.url="https://ghcr.io/lao-tseu-is-alive/go-cloud-k8s-jwt-login:latest"
LABEL org.opencontainers.image.authors="cgil"
LABEL org.opencontainers.image.licenses="MIT"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY cmd/server ./server
COPY pkg ./pkg


# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o go-cloud-k8s-jwt-login-server ./server


######## Start a new stage  #######
FROM scratch
USER 12121:12121
WORKDIR /goapp
# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/go-cloud-k8s-jwt-login-server .

# Expose port to the outside world, server will use the env PORT as listening port or 8888 as default
EXPOSE 8888

# Health check to ensure the app is running
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8888/health || exit 1

# Command to run the executable
CMD ["./go-cloud-k8s-jwt-login-server"]
