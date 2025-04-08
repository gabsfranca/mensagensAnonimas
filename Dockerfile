FROM golang:1.23.7-alpine AS builder

WORKDIR /app

COPY backend/go.mod backend/go.sum ./

RUN go mod download && go mod tidy

COPY backend/ ./

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .

COPY frontend/ ./frontend/

RUN touch .env

EXPOSE 8080

CMD ["./main"]