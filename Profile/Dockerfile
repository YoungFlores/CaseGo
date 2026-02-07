FROM golang:alpine AS builder
WORKDIR /app
ARG SERVICE_NAME
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/${SERVICE_NAME}/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations
CMD ["./main"]