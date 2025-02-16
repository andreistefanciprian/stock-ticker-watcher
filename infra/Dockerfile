# Build stage
FROM golang:1.24-alpine AS build

WORKDIR /app

COPY go.* ./
RUN go mod download
COPY main.go .

ENV GOOS=linux \
    GOARCH=amd64

RUN go build -a -o stock-ticker-watcher

# Deploy stage
FROM gcr.io/distroless/base-debian12

COPY --from=build /app/stock-ticker-watcher .

USER nonroot:nonroot

ENTRYPOINT ["./stock-ticker-watcher"]