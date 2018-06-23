FROM golang:1.10.0-alpine3.7
MAINTAINER Grant Ellis <robert.grant.ellis@gmail.com>

ENV appName t9
ENV packageName github.com/RobertGrantEllis/${appName}
ENV storagePath /data

ADD . ${GOPATH}/src/${packageName}

RUN go build -o /${appName} ${packageName} && rm -rf /usr/local/bin/* /usr/local/go ${GOPATH}

EXPOSE 4239

CMD /${appName} server --address=0.0.0.0:4239