FROM golang:1.19-alpine AS builder

WORKDIR /mux-rpc

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o mux-rpc-log-server ./cmd/server/main.go
RUN go build -o mux-rpc-log-client ./client/client.go

#####################################
FROM alpine:3.17

WORKDIR /mux-rpc

COPY --from=builder /mux-rpc/mux-rpc-log-server .
COPY --from=builder /mux-rpc/mux-rpc-log-client .
COPY wait.sh . 

RUN chmod +x wait.sh mux-rpc-log-client

EXPOSE 8080
EXPOSE 3214

CMD ["./mux-rpc-log-server"]
