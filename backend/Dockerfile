FROM golang:1.23.7-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/cmd

RUN go build -o /app/app .

CMD ["/app/app"]
