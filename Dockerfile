FROM node:22-alpine AS build

WORKDIR /app
COPY hugo/package.json hugo/package-lock.json ./
RUN npm ci --omit=dev --omit optional

ARG DO_NOT_MINIFY

ENV NODE_ENV=production

COPY ./hugo/webpack.mix.js ./hugo/tsconfig.json ./hugo/.babelrc.js /app/
COPY ./hugo/src/ /app/src/
COPY ./hugo/layouts /app/layouts/

RUN npm run build

FROM golang:1.24-alpine AS go-build
COPY rss_generator /build
RUN cd /build && go build -o /build/bin/rss_generator -ldflags "-s -w" && ls -la /build/bin/rss_generator

FROM alpine:3.19

# enables automatic changelog generation by tools like Dependabot
LABEL org.opencontainers.image.source="https://github.com/radio-t/radio-t-site"

RUN \
    apk add --update --no-cache tzdata curl openssl git openssh-client ca-certificates && \
    cp /usr/share/zoneinfo/EST /etc/localtime && \
    echo "CDT" > /etc/timezone && date && \
    rm -rf /var/cache/apk/*

ENV HUGO_VER=0.81.0
ADD https://github.com/gohugoio/hugo/releases/download/v${HUGO_VER}/hugo_${HUGO_VER}_Linux-64bit.tar.gz /tmp/hugo.tar.gz
RUN \
    cd /tmp && tar -zxf hugo.tar.gz && ls -la && \
    cp -fv /tmp/hugo /bin/hugo

COPY --from=build /app/static/build/ /app/static/build/
COPY --from=build /app/data/manifest.json /app/data/manifest.json

COPY --from=go-build /build/bin/rss_generator /usr/local/bin/rss_generator

COPY exec.sh /srv/exec.sh
RUN chmod +x /srv/exec.sh

CMD ["/srv/exec.sh"]
