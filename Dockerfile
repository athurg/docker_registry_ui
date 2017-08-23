FROM golang AS compile
ADD . $GOPATH/src/app
RUN go get app
RUN CGO_ENABLED=0 go install -a app

FROM scratch
ENV DSN ""
ENV LISTEN_ADDR ":80"
ENV TOKEN_ISSUER "AuthServer"
ENV KEY_PEM_BLOCK ""
ENV CERT_PEM_BLOCK ""
ENV USER_TABLE_NAME "users"
ENV PRIVILEGE_TABLE_NAME "privileges"
ENV TOKEN_EXPIRATION "900"
ENV DSN "user:pass@tcp(host:port)/dbname?charset=utf8&parseTime=True&loc=Local"

COPY --from=compile /go/bin/app /app
CMD ["/app"]
