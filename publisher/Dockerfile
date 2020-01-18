FROM python:3.8-alpine

ENV \
    TERM=xterm-color           \
    TIME_ZONE=America/Chicago  \
    MYUSER=app                 \
    MYUID=1001

RUN \
    apk add --no-cache --update su-exec tzdata curl openssl git lftp openssh && \
    ln -s /sbin/su-exec /usr/local/bin/gosu && \
    mkdir -p /home/$MYUSER && \
    adduser -s /bin/sh -D -u $MYUID $MYUSER && chown -R $MYUSER:$MYUSER /home/$MYUSER && \
    delgroup ping && addgroup -g 998 ping && \
    mkdir -p /srv && chown -R $MYUSER:$MYUSER /srv && \
    cp /usr/share/zoneinfo/${TIME_ZONE} /etc/localtime && \
    echo "${TIME_ZONE}" > /etc/timezone && date && \
    rm -rf /var/cache/apk/*

COPY . /srv/publisher
WORKDIR /srv/publisher

RUN pip install --no-cache-dir -r requirements.txt

USER $MYUSER

ENTRYPOINT ["invoke"]
CMD ["--list"]
