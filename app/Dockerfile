FROM golang:1.21

RUN mkdir /dq
ADD . /dq
WORKDIR /dq

RUN go get -d ./...

RUN go build -o dq .

ENTRYPOINT /dq/dq
EXPOSE 12000
