FROM golang:1.8
LABEL maintainer "Kapeel Sable <kapeel.sable@gmail.com>"

RUN echo 'deb http://ftp.de.debian.org/debian jessie main' >> /etc/apt/sources.list \
    && apt-get update \
    && apt-get install libpam0g-dev

VOLUME /go/src/pam_ohmage
COPY ./docker_entrypoint.sh /docker_entrypoint.sh
RUN chmod +x /docker_entrypoint.sh

WORKDIR /go/src/pam_ohmage
ENTRYPOINT /docker_entrypoint.sh