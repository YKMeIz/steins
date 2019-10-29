FROM golang:alpine

RUN apk add --no-cache tzdata

WORKDIR /
COPY ./steins .
COPY ./server.crt .
COPY ./server.key .

EXPOSE 443/tcp


ENTRYPOINT ["/steins"]
