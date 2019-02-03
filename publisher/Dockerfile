FROM golang:latest AS build

WORKDIR /go/src/github.com/radio-t/radio-t-site/publisher

RUN go get -u google.golang.org/api/youtube/v3 && \
    go get -u golang.org/x/oauth2 && \
    go get -u golang.org/x/oauth2/google && \
    go get -u golang.org/x/net/context && \
    go get -u google.golang.org/api/compute/v1 && \
    go get -u github.com/spf13/cobra/cobra

COPY podcast-to-youtube podcast-to-youtube
COPY add-to-youtube add-to-youtube

RUN go install github.com/radio-t/radio-t-site/publisher/add-to-youtube

FROM jrottenberg/ffmpeg:latest

WORKDIR /publisher

COPY --from=build /go/bin/add-to-youtube /usr/local/bin/
COPY entrypoint.sh .

ENTRYPOINT [ "/publisher/entrypoint.sh" ]