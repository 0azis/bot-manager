FROM golang:1.21.9

COPY . /app
WORKDIR /app/cmd/

RUN go mod download

RUN go build -o ./server

EXPOSE 3000

CMD ["./server"]