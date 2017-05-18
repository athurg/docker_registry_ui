#Dockerfile for build
FROM golang
MAINTAINER Feng Jianbo <fengjianbo@nibirutech.com>

ADD . $GOPATH/src/app
RUN go get app
RUN CGO_ENABLED=0 go install -a app
