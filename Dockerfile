FROM golang:1.14

WORKDIR /go/src/juntosajugar
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

RUN go run ./cmd/web
