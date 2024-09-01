FROM golang:1.20-buster

WORKDIR /app

ENV GOOS=darwin
ENV GOARCH=amd64

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o commit-log ./cmd

CMD ["./commit-log"]
