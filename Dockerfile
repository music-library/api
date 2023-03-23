# build environment
FROM golang:1.18.3-alpine as build
WORKDIR /app
COPY / /app
RUN apk add --update gcc libc-dev git
RUN go install -mod vendor github.com/go-task/task/v3/cmd/task
RUN go generate -tags tools tools/tools.go
RUN task bootstrap
RUN task build

# production environment
FROM alpine:3.16.0
RUN apk add vips-tools mediainfo
COPY --from=build /app/bin /app
WORKDIR /app
EXPOSE 80
ENV PORT 80
ENV HOST 0.0.0.0
ENV DATA_DIR "./data"
ENV MUSIC_DIR "./music"
VOLUME ["./data"]
VOLUME ["./music"]
CMD ["./music-api"]
