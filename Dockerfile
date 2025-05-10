FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go mod tidy

RUN go build -o main-web ./cmd/web/main.go
# RUN go build -o main-worker ./cmd/worker/main.go  # Uncomment jika ada worker

RUN mkdir -p ./storage/logs ./temp && chmod -R 777 ./storage/logs ./temp

COPY temp temp
COPY app.env app.env

EXPOSE 3000