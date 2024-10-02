FROM golang:1.20

WORKDIR /app

RUN apt-get update && apt-get install -y librdkafka-dev \
    && rm -rf /var/lib/apt/lists/*  

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV DB_USERNAME="root"
ENV DB_PASSWORD="root"
ENV DB_HOST="wallet_core_db"
ENV DB_PORT="3306"
ENV DB_NAME="wallet_core_db"


CMD ["go", "run", "./cmd/wallet_balance/main.go"]

