FROM node:10-alpine as build

WORKDIR /app

COPY hugo/package.json hugo/package-lock.json ./
RUN npm ci

ENV NODE_ENV=production
COPY hugo/.modernizr.js hugo/webpack.mix.js hugo/tsconfig.json hugo/.babelrc.js ./
COPY hugo/src/ src/
RUN npm run production

###
###

FROM alpine:3.8

RUN \
    apk add --update --no-cache tzdata curl openssl git openssh-client python3 ca-certificates && \
    apk add --no-cache --virtual .build-deps python3-dev && \
    python3 -m ensurepip && pip3 install --upgrade pip && \
    pip3 install pytoml mistune && \
    apk del .build-deps && \
    cp /usr/share/zoneinfo/EST /etc/localtime && \
    echo "CDT" > /etc/timezone && date && \
    rm -rf /var/cache/apk/*

ENV HUGO_VER=0.80.0
ADD https://github.com/gohugoio/hugo/releases/download/v${HUGO_VER}/hugo_${HUGO_VER}_Linux-64bit.tar.gz /tmp/hugo.tar.gz
RUN \
    cd /tmp && tar -zxf hugo.tar.gz && ls -la && \
    cp -fv /tmp/hugo /bin/hugo

COPY --from=build /app/static/build/ /app/static/build/
COPY --from=build /app/data/manifest.json /app/data/manifest.json

## COPY --chmod=0755 exec.sh /srv/exec.sh worked for docker CE >= 20.10
## Implemented at https://github.com/moby/buildkit/pull/1492
## Current version of docker daemon in Github Actions is 19.03.13+azure
COPY --chmod=0755 exec.sh /srv/exec.sh
# COPY exec.sh /srv/exec.sh
# RUN chmod +x /srv/exec.sh

CMD ["/srv/exec.sh"]
