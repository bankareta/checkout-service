FROM golang:1.23.4

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN apt-get update && apt-get install -y netcat-openbsd

RUN chmod +x /app/wait-for.sh

RUN go mod tidy

RUN go build -o main-web ./cmd/web/main.go

RUN mkdir -p ./storage/logs ./temp && chmod -R 777 ./storage/logs ./temp

COPY app.env app.env

EXPOSE 3001