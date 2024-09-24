FROM golang:1.20

WORKDIR /app

RUN apt-get update && apt-get install -y librdkafka-dev \
    && rm -rf /var/lib/apt/lists/* 

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/wallet_core ./cmd/wallet_core/main.go

RUN chmod +x /app/wallet_core

ENV DB_USERNAME="root"
ENV DB_PASSWORD="root"
ENV DB_HOST="localhost"
ENV DB_PORT="3306"
ENV DB_NAME="wallet"

EXPOSE 8080

CMD ["/app/wallet_core"]
