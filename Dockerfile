FROM ubuntu:22.04

WORKDIR /usr/bin

RUN mkdir /var/log/bot

COPY <PATH_TO_BOT_BINARY_FILE> .

RUN apt-get update && apt-get install -y ca-certificates && update-ca-certificates

ENTRYPOINT ["<PATH_TO_BOT_BINARY_FILE>"]