# syntax=docker/dockerfile:1

FROM golang:1.21.6 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
# run the build
RUN CGO_ENABLED=0 GOOS=linux go build -o /ginrcon



# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test -v ./...

# Deploy the application binary into a lean image
FROM alpine:3.14 AS build-release-stage

WORKDIR /

COPY --from=build-stage /ginrcon /ginrcon

EXPOSE 8080
ENTRYPOINT ["/ginrcon"]
