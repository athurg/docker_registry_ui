FROM golang AS compile
ADD . $GOPATH/src/app
RUN go get app
RUN CGO_ENABLED=0 go install -a app

FROM registry:2
COPY --from=compile /go/bin/app /ui
RUN sed -i.bak '1a/ui &' /entrypoint.sh
