FROM golang AS compile
ADD . $GOPATH/src/app
RUN go get app
RUN CGO_ENABLED=0 go install -a app

FROM scratch
COPY --from=compile /go/bin/app /app
CMD ["/app"]
