FROM golang:1.9.0-alpine3.6 as build

RUN apk add --no-cache --virtual git musl-dev
RUN go get github.com/golang/dep/cmd/dep

RUN go get github.com/gohugoio/hugo
WORKDIR /go/src/github.com/gohugoio/hugo
RUN dep ensure
RUN go install -ldflags '-s -w'

FROM alpine:3.6

COPY --from=build /go/bin/hugo /bin/hugo
COPY exec.sh /srv/exec.sh

RUN \
 chmod +x /srv/exec.sh && \
 apk add --update --no-cache tzdata curl openssl git openssh-client python3 ca-certificates && \
 apk add --no-cache --virtual .build-deps python3-dev && \
 python3 -m ensurepip && pip3 install --upgrade pip && \
 pip3 install pytoml mistune && \
 apk del .build-deps && \
 cp /usr/share/zoneinfo/EST /etc/localtime && \
 echo "CDT" > /etc/timezone && date && \
 rm -rf /var/cache/apk/*

CMD ["/srv/exec.sh"]
