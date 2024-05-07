# syntax=docker/dockerfile:1

FROM golang:1.21

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go mod tidy
RUN go build -C services/worker/service -o myapp

EXPOSE 8081 8082 8083

CMD ["services/worker/service/myapp"]