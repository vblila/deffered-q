FROM golang:1.21

RUN mkdir /dq
ADD . /dq
WORKDIR /dq

RUN go get -d ./...

RUN go build -o dq .

RUN mkdir /var/log/dq

ENTRYPOINT /dq/dq \
    -h "" \
    -p 12000 \
    -debug ${DQ_PROFILER_ENABLED} \
    -ict ${DQ_INACTIVE_CONNECTION_TIME_SECONDS} \
    -rtt ${DQ_RESERVED_TASK_STUCK_TIME_SECONDS} \
    -rta ${DQ_RESERVED_TASK_STUCK_MAX_ATTEMPTS} \
    >> /var/log/dq/server.log 2>&1

EXPOSE 12000
