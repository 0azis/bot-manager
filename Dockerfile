FROM golang:1.21.9

COPY . /app
WORKDIR /app/cmd/

RUN go mod download

RUN go build -o ./server

EXPOSE 8000 

CMD ["./server"]
