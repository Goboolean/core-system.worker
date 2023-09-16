FROM golang

WORKDIR /go/src/app

ADD . .

RUN go mod tidy

RUN go build -o main cmd/main/run.go

CMD ["./main"] 