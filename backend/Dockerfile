FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /franchise-api ./cmd/api

FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /franchise-api .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/.env .env

RUN mkdir -p /app/uploads

EXPOSE 8080

CMD ["./franchise-api"]
