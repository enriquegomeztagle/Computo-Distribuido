FROM golang:1.21

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN apt-get update && apt-get install -y golang-cfssl

RUN make init

RUN make gencert

RUN chmod +x run_tests.sh

CMD ["./run_tests.sh"]
