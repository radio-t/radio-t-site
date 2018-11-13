FROM alpine:3.8

ENV \
    TERM=xterm-color           \
    TIME_ZONE=America/Chicago  \
    MYUSER=app                 \
    MYUID=1000

RUN \
    apk add --no-cache --update su-exec tzdata curl openssl && \
    ln -s /sbin/su-exec /usr/local/bin/gosu && \
    mkdir -p /home/$MYUSER && \
    adduser -s /bin/sh -D -u $MYUID $MYUSER && chown -R $MYUSER:$MYUSER /home/$MYUSER && \
    delgroup ping && addgroup -g 998 ping && \
    mkdir -p /srv && chown -R $MYUSER:$MYUSER /srv && \
    cp /usr/share/zoneinfo/${TIME_ZONE} /etc/localtime && \
    echo "${TIME_ZONE}" > /etc/timezone && date && \
    rm -rf /var/cache/apk/*

ADD updater.sh /srv/updater.sh
ADD entrypoint.sh /srv/entrypoint.sh

RUN \
    apk add --update openssh-client && \
    mkdir -p /home/app/.ssh && \
    echo "StrictHostKeyChecking=no" > /home/app/.ssh/config && \
    chown -R app:app /home/app/.ssh/ && \
    chmod 600 /home/app/.ssh/* && \
    chmod 700 /home/app/.ssh  && \
    chmod +x /srv/updater.sh /srv/entrypoint.sh

USER app
CMD ["/srv/entrypoint.sh"]
