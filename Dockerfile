FROM golang:alpine

ARG pkg=github.com/pubkraal/ostls

RUN apk add --no-cache ca-certificates

ENV GO111MODULE=on

COPY . $GOPATH/src/$pkg

RUN set -ex \
      && apk add --no-cache --virtual .build-deps \
              git \
      && go build $pkg \
      && go install $pkg \
      && apk del .build-deps

CMD echo "Use `ostls'"; exit 1
