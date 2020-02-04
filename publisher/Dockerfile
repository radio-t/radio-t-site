FROM umputun/baseimage:buildgo-latest as build

ADD app /build
WORKDIR  /build
RUN go test -mod=vendor ./...
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o /target/publisher -ldflags \
    "-X main.revision=$(git rev-parse --abbrev-ref HEAD)-$(git describe --abbrev=7 --always --tags)-$(date +%Y%m%d-%H:%M:%S)"


FROM umputun/baseimage:app-latest

RUN \
    apk add --no-cache --update python3 py-pip git && \
    rm -rf /var/cache/apk/* && \
    pip3 install --upgrade pip

RUN \
    apk add --update openssh-client && \
    mkdir -p /home/app/.ssh && \
    echo "StrictHostKeyChecking=no" > /home/app/.ssh/config && \
    chown -R app:app /home/app/.ssh/ && \
    chmod 600 /home/app/.ssh/* && \
    chmod 700 /home/app/.ssh

COPY --from=build /target/publisher /srv/publisher/publisher
COPY . /srv/publisher
WORKDIR /srv/publisher

RUN pip3 install --no-cache-dir -r requirements.txt

USER app

ENTRYPOINT ["invoke"]
CMD ["--list"]
