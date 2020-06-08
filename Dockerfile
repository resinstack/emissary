FROM golang:alpine3.12 as build
RUN mkdir -p /go/emissary && \
        apk add upx
COPY ./ /go/emissary
RUN cd /go/emissary && \
        go mod vendor && \
        CGO_ENABLED=0 go build -x -ldflags '-s -w' -o /emissary . && \
        ls -alh /emissary && \
        upx /emissary && \
        ls -alh /emissary

FROM alpine:latest as certs
RUN apk --update add ca-certificates

FROM scratch
COPY --from=build /emissary /
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
CMD ["/emissary"]
