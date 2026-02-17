FROM golang:1.22 AS builder
ENV CGO_ENABLED=1 GOOS=linux
RUN apt-get update && apt-get install -y gcc libc6-dev libsqlite3-dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -ldflags '-linkmode external -extldflags "-static"' -o /app/main cmd/gatemulator/main.go

FROM scratch
WORKDIR /app
COPY --from=builder /app/main /app/main
EXPOSE 34000
CMD ["./main"]
