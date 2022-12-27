# build environment
FROM golang:1.18.3-alpine as build
WORKDIR /app
COPY / /app
RUN apk add --update make gcc g++ libc-dev
ENV PKG_CONFIG_PATH /usr/lib:/usr/local/lib/pkgconfig:/usr/lib/pkgconfig:$PKG_CONFIG_PATH
RUN make bootstrap && make build

# production environment
FROM alpine:3.16.0
RUN apk add vips-dev
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
