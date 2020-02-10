FROM umputun/baseimage:buildgo-latest as build

ARG REV
ADD app /build
WORKDIR  /build

RUN go test -mod=vendor ./...
RUN echo "version=${REV}"
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o /target/publisher -ldflags "-X main.revision=${REV}"


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

COPY --from=build /target/publisher /usr/local/bin/publisher
COPY . /srv/publisher
RUN chmod 600 /srv/publisher/mp3tags/cmd.py && \
    ln -s /srv/publisher/mp3tags/cmd.py /usr/local/bin/mp3tags

WORKDIR /srv/hugo

RUN pip3 install --no-cache-dir -r /srv/publisher/mp3tags/requirements.txt

USER app

ENTRYPOINT ["/usr/local/bin/publisher"]
