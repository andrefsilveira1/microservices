FROM golang:1.20

WORKDIR /app

RUN apt-get update && apt-get install -y librdkafka-dev \
    && rm -rf /var/lib/apt/lists/*  

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/wallet_balance ./cmd/wallet_balance/main.go

EXPOSE 8000

ENV MYSQL_ROOT_PASSWORD="root"
ENV MYSQL_DATABASE="wallet_balance_db"
ENV MYSQL_USER="root"
ENV MYSQL_PASSWORD="root"

CMD ["/app/wallet_balance"]