FROM golang:alpine as BUILD

WORKDIR /go/src/github.com/randomcoww/go-mpd-es
COPY . .

RUN set -x \
  \
  && apk add --no-cache \
    git \
  \
  && go get -d ./... \
  && go build

FROM alpine:edge

COPY --from=BUILD /go/src/github.com/randomcoww/go-mpd-es/go-mpd-es /

ENTRYPOINT ["/go-mpd-es"]
