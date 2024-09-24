FROM golang:1.20

WORKDIR /app

RUN apt-get update && apt-get install -y librdkafka-dev \
    && rm -rf /var/lib/apt/lists/*  

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/wallet_balance ./cmd/wallet_balance/main.go

EXPOSE 8000

ENV DB_USERNAME="root"
ENV DB_PASSWORD="root"
ENV DB_HOST="localhost"
ENV DB_PORT="3307"
ENV DB_NAME="wallet_balance_db"

CMD ["/app/wallet_balance"]