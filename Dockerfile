# syntax=docker/dockerfile:1

# run the build
FROM golang:1.21.6 AS build-stage
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /ginrcon

# run tests
FROM build-stage AS run-test-stage
RUN go test -v ./...

# prep runtime
FROM alpine:3.14 AS build-release-stage
WORKDIR /
COPY --from=build-stage /ginrcon /ginrcon
ENV PORT=8080 \
    TRUSTED_PROXIES= \
    RCON_SERVER= \
    RCON_PORT= \
    RCON_ADMIN_PASSWORD=

EXPOSE 8080
ENTRYPOINT ["/ginrcon"]
