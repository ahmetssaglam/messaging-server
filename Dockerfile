FROM golang:1.24 AS builder

WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends git \
    && rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd

FROM alpine:3.18

WORKDIR /app

# install CA certs if your app makes HTTPS calls
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/server .

EXPOSE 8080

ENTRYPOINT ["./server"]
