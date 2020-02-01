FROM umputun/baseimage:app-latest

RUN \
    apk add --no-cache --update python3 py-pip git openssh && \
    rm -rf /var/cache/apk/* && \
    pip3 install --upgrade pip

COPY . /srv/publisher
WORKDIR /srv/publisher

RUN pip3 install --no-cache-dir -r requirements.txt

USER app

ENTRYPOINT ["invoke"]
CMD ["--list"]
