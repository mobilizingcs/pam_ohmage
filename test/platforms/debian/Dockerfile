FROM debian:jessie
LABEL maintainer "Kapeel Sable <kapeel.sable@gmail.com>"

ENV RSTUDIO_PRO_VERSION 1.0.143

RUN set -x \
    && export DEBIAN_FRONTEND=noninteractive \
    && echo 'deb http://cran.rstudio.com/bin/linux/debian jessie-cran3/' >> /etc/apt/sources.list \
    && apt-key adv --keyserver keys.gnupg.net --recv-key 381BA480 \
    && apt-get update \
    && apt-get install -y r-base gdebi libcurl4-openssl-dev libssl-dev locales syslog-ng

RUN echo "en_US.UTF-8 UTF-8" > /etc/locale.gen && \
    locale-gen en_US.UTF-8 && \
    /usr/sbin/update-locale LANG=en_US.UTF-8
ENV LC_ALL en_US.UTF-8

RUN curl -fSL https://download2.rstudio.org/rstudio-server-pro-${RSTUDIO_PRO_VERSION}-amd64.deb -o /tmp/rstudio-server.deb \
    && gdebi -n /tmp/rstudio-server.deb \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

COPY ./pam-d-samples/rstudio /etc/pam.d/rstudio
COPY ./pam-d-samples/rstudio-session /etc/pam.d/rstudio-session
RUN echo "auth-pam-sessions-profile=rstudio-session" >> /etc/rstudio/rserver.conf

VOLUME /pam_ohmage/bin

EXPOSE 8787
COPY ./test/platforms/test_entrypoint.sh /test_entrypoint.sh
RUN chmod +x /test_entrypoint.sh
RUN ln -s /pam_ohmage/bin/pam_ohmage.so /lib/security/pam_ohmage.so

RUN sed -i -E 's/^(\s*)system\(\);/\1unix-stream("\/dev\/log");/' /etc/syslog-ng/syslog-ng.conf
RUN sed -i 's/^#\(SYSLOGNG_OPTS="--no-caps"\)/\1/g' /etc/default/syslog-ng

ENTRYPOINT /test_entrypoint.sh