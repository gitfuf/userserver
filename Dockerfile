FROM golang:1.10.3-alpine3.7
MAINTAINER Fuf

# Install tools required for project
RUN apk update && apk upgrade && apk add --no-cache bash git
RUN go get -u github.com/gorilla/mux \
	&& go get -u github.com/denisenkom/go-mssqldb \
	&& go get -u github.com/go-sql-driver/mysql \
	&& go get -u github.com/lib/pq \
	&& go get -u github.com/spf13/viper

ENV SOURCES /go/src/github.com/gitfuf/userserver/
COPY . ${SOURCES}

RUN cd ${SOURCES}cmd/server && CGO_ENABLED=0 go build

WORKDIR ${SOURCES}cmd/server
CMD ${SOURCES}cmd/server/server

