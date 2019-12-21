FROM golang:1.13-alpine3.10
RUN apk add git
WORKDIR /go/src/test

COPY . .

#RUN go get -d -v ./...
RUN go get -t -v ./...
#RUN go install -v ./...
ENV CGO_ENABLED 0

CMD ["./run"]
