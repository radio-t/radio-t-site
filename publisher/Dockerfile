FROM golang:latest AS build

WORKDIR /publisher

COPY add-to-youtube add-to-youtube

WORKDIR /publisher/add-to-youtube

RUN go get ./... && go install

FROM jrottenberg/ffmpeg:latest

WORKDIR /publisher

COPY --from=build /go/bin/add-to-youtube /usr/local/bin/
COPY entrypoint.sh .

ENTRYPOINT [ "/publisher/entrypoint.sh" ]