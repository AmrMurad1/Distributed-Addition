FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o service-a .

FROM alpine:latest

RUN apk --no-cache add ca-certificates && \
    adduser -D -s /bin/sh appuser

WORKDIR /app
COPY --from=builder /app/service-a .
RUN chown appuser:appuser service-a

USER appuser
EXPOSE 50051

CMD ["./service-a"]
