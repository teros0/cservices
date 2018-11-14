FROM golang

ADD . /go/src/cservices/

RUN go install cservices/

ENTRYPOINT /go/bin/cservices -path=/data.csv -port=:7777

COPY ./resources/data.csv /

EXPOSE 7777