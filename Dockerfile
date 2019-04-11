
FROM alpine

RUN apk update && \
    mkdir -p /run/docker/plugins

COPY l2vlan l2vlan