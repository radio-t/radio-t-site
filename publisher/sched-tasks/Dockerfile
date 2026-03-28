FROM umputun/baseimage:buildgo-latest AS build

WORKDIR /build
COPY app ./app
COPY go.mod go.sum ./
COPY vendor ./vendor

RUN \
    cd /build && \
    go build -o /build/bin/publisher -ldflags "-s -w" ./app/main.go

FROM umputun/cronn:latest

USER root

RUN \
    apk add --update openssh-client git bash && \
    mkdir -p /home/app/.ssh && \
    echo "StrictHostKeyChecking=no" > /home/app/.ssh/config && \
    chown -R app:app /home/app/.ssh/ && \
    chmod 600 /home/app/.ssh/config && \
    chmod 700 /home/app/.ssh

COPY --from=build /build/bin/publisher /usr/local/bin/publisher

USER app
WORKDIR /srv

RUN git config --global --add safe.directory /srv