FROM golang:1.11

WORKDIR /go/src/teste
COPY . .

RUN go get github.com/gomodule/redigo/redis

RUN go get github.com/gorilla/mux

RUN go install -v ./...

EXPOSE 8080

CMD ["teste"]
