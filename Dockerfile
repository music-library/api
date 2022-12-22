# build environment
FROM golang:1.18.3-alpine as build
WORKDIR /app
COPY / /app
RUN apk add --update make gcc libc-dev
RUN make bootstrap && make build

# production environment
FROM alpine:3.16.0
COPY --from=build /app/bin /app
WORKDIR /app
EXPOSE 80
ENV PORT 80
ENV HOST 0.0.0.0
CMD ["./music-api"]
