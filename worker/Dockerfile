FROM golang
WORKDIR /go/src/github.com/clitetailor/distributed-hash-decrypter
COPY . .
WORKDIR /go/src/github.com/clitetailor/distributed-hash-decrypter/worker
RUN go build
CMD ["./worker", "--host", "master"]
