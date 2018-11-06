FROM golang

ADD . /go/src/github.com/teros0/cservices/

RUN go install github.com/teros0/cservices/

ENTRYPOINT ["/go/bin/cservices"]

EXPOSE 7777