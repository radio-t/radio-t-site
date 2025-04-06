FROM umputun/baseimage:buildgo-latest AS build
ARG GIT_BRANCH
ARG GITHUB_SHA
ARG CI
ARG REV

ADD . /build

RUN \
    cd /build/app &&\
    go test ./...

RUN \
    if [ -z "$CI" ] ; then echo "runs outside of CI" && version=${REV}; \
    else version=${GIT_BRANCH}-${GITHUB_SHA:0:7}-$(date +%Y%m%dT%H:%M:%S); fi && \
    echo "version=$version" && \
    cd  /build/app && go build -o /build/bin/publisher -ldflags "-X main.revision=${version} -s -w" &&\
    ls -la /build/bin/publisher


FROM umputun/baseimage:app-latest

# enables automatic changelog generation by tools like Dependabot
LABEL org.opencontainers.image.source="https://github.com/radio-t/radio-t-site"

RUN \
    apk add --update openssh-client git && \
    mkdir -p /home/app/.ssh && \
    echo "StrictHostKeyChecking=no" > /home/app/.ssh/config && \
    chown -R app:app /home/app/.ssh/ && \
    chmod 600 /home/app/.ssh/* && \
    chmod 700 /home/app/.ssh

RUN \
    export SPOT_VERSION=v1.16.4 && \
    wget https://github.com/umputun/spot/releases/download/${SPOT_VERSION}/spot_${SPOT_VERSION}_linux_amd64.apk && \
    ls -la spot_${SPOT_VERSION}_linux_amd64.apk && \
    apk add --allow-untrusted  spot_${SPOT_VERSION}_linux_amd64.apk && \
    rm spot_${SPOT_VERSION}_linux_amd64.apk && \
    spot --version

COPY --from=build /build/bin/publisher /usr/local/bin/publisher
COPY spot.yml /etc/spot.yml

WORKDIR /srv/hugo
USER app

RUN git config --global --add safe.directory /srv

CMD ["-h"]
ENTRYPOINT ["/usr/local/bin/publisher"]
