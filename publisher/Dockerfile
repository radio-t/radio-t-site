FROM umputun/baseimage:buildgo-latest as build
ARG GIT_BRANCH
ARG GITHUB_SHA
ARG CI
ARG REV
ENV GOFLAGS="-mod=vendor"

ADD . /build
#
#RUN \
#    cd /build/app &&\
#    go test ./...

RUN \
    if [ -z "$CI" ] ; then echo "runs outside of CI" && version=${REV}; \
    else version=${GIT_BRANCH}-${GITHUB_SHA:0:7}-$(date +%Y%m%dT%H:%M:%S); fi && \
    echo "version=$version" && \
    cd  /build/app && go build -o /build/bin/publisher -ldflags "-X main.revision=${version} -s -w" &&\
    ls -la /build/bin/publisher


FROM umputun/baseimage:app-latest

RUN \
    apk add --update openssh-client git && \
    mkdir -p /home/app/.ssh && \
    echo "StrictHostKeyChecking=no" > /home/app/.ssh/config && \
    chown -R app:app /home/app/.ssh/ && \
    chmod 600 /home/app/.ssh/* && \
    chmod 700 /home/app/.ssh

RUN \
    wget https://github.com/umputun/spot/releases/download/v1.12.0/spot_v1.12.0_linux_amd64.apk && \
    ls -la spot_v1.12.0_linux_amd64.apk && \
    apk add --allow-untrusted  spot_v1.12.0_linux_amd64.apk && \
    rm spot_v1.12.0_linux_amd64.apk && \
    spot --version

COPY --from=build /build/bin/publisher /usr/local/bin/publisher
COPY spot.yml /etc/spot.yml

WORKDIR /srv/hugo
USER app

RUN git config --global --add safe.directory /srv

CMD ["-h"]
ENTRYPOINT ["/usr/local/bin/publisher"]
