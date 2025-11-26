FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/api ./cmd/api/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/worker ./cmd/worker/main.go

FROM alpine:latest AS api

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/api .

RUN mkdir -p /app/storage/photos

EXPOSE 8080

CMD ["./api"]

FROM alpine:latest AS worker

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/worker .

CMD ["./worker"]
