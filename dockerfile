FROM alpine:latest

MAINTAINER YuCheng

USER root

ADD run /
ADD config.yml /

EXPOSE 9100

ENTRYPOINT ["/run"]