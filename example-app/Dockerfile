FROM golang:1.13-alpine3.10
RUN apk add git
WORKDIR /go/src/app
ADD metrics.json /var/cadvisor/metrics.json
LABEL io.cadvisor.metric.prometheus="/var/cadvisor/metrics.json"

COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

ENTRYPOINT ["app"]
#CMD ["-c", "40"]
