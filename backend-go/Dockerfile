FROM golang:1.25-alpine AS deps
RUN apk add --no-cache git ca-certificates tzdata
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

FROM golang:1.25-alpine AS builder
RUN apk add --no-cache git ca-certificates tzdata
WORKDIR /app
COPY --from=deps /go/pkg/mod /go/pkg/mod
COPY go.mod go.sum ./
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main ./cmd

FROM alpine:latest AS runtime
RUN apk --no-cache add ca-certificates curl bash wget && \
    adduser -D -s /bin/sh appuser
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/templates ./templates
RUN mkdir -p uploads && chown -R appuser:appuser /app
USER appuser
EXPOSE 5000
CMD ["./main"]
