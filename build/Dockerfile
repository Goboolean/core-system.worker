FROM golang

WORKDIR /go/src/app

ADD . .

RUN go mod tIDy

RUN go build -o main cmd/main/run.go

CMD ["./main"] 