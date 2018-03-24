FROM golang
WORKDIR /go/src/github.com/clitetailor/distributed-hash-decrypter
COPY . .
WORKDIR /go/src/github.com/clitetailor/distributed-hash-decrypter/master
RUN go build
CMD ["./master"]
