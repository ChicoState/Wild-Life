FROM golang:1.18

WORKDIR /go/src
COPY . .
RUN apt-get install pkg-config
RUN go get -u -d gocv.io/x/gocv
RUN go get -d -v ./...
RUN go install ./main.go

CMD ["main"]