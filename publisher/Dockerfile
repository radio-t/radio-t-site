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

RUN git config --global --add safe.directory /srv

COPY . /srv/publisher
WORKDIR /srv/publisher

RUN pip3 install --no-cache-dir -r requirements.txt

USER app

ENTRYPOINT ["invoke"]
CMD ["--list"]
