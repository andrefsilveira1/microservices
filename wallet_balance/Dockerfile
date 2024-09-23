FROM golang:1.20

WORKDIR /app

RUN apt-get update && apt-get install -y librdkafka-dev \
    && rm -rf /var/lib/apt/lists/*  

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o wallet_balance ./cmd/main.go

EXPOSE 8000

CMD ["./wallet_balance"]
