FROM golang:1.21

RUN mkdir /dq
ADD . /dq
WORKDIR /dq

RUN go get -d ./...

RUN go build -o dq .

RUN mkdir /var/log/dq

ENTRYPOINT /dq/dq >> /var/log/dq/server.log 2>&1

EXPOSE 12000
