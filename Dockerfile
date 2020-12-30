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

FROM alpine:3.12
RUN apk add --no-cache ca-certificates
COPY --from=build /emissary /
CMD ["/emissary"]
