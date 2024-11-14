# syntax=docker/dockerfile:1

# run the build
FROM golang:1.22.9 AS build-stage
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /ginrcon

# run tests
FROM build-stage AS run-test-stage
RUN go test -v ./...

# prep runtime
FROM alpine:3.20 AS build-release-stage
WORKDIR /
COPY --from=build-stage /ginrcon /ginrcon
ENV PORT=8080

EXPOSE 8080/tcp
ENTRYPOINT ["/ginrcon"]
