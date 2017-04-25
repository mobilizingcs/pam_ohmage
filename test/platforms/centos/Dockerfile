FROM centos:7
LABEL maintainer "Kapeel Sable <kapeel.sable@gmail.com>"

ENV RSTUDIO_PRO_VERSION 1.0.143

RUN set -x \
    && yum -y update \
    && yum -y install epel-release \
    && yum -y install R syslog-ng

#RUN localectl set-locale LANG=en_US.UTF-8
ENV LC_ALL en_US.UTF-8

RUN curl -fSL https://download2.rstudio.org/rstudio-server-rhel-pro-{$RSTUDIO_PRO_VERSION}-x86_64.rpm -o /tmp/rstudio-server.rpm
RUN yum -y install --nogpgcheck /tmp/rstudio-server.rpm

COPY ./pam-d-samples/rstudio /etc/pam.d/rstudio
COPY ./pam-d-samples/rstudio-session /etc/pam.d/rstudio-session
RUN echo "auth-pam-sessions-profile=rstudio-session" >> /etc/rstudio/rserver.conf

VOLUME /pam_ohmage/bin

EXPOSE 8787
COPY ./test/platforms/centos/docker_entrypoint.sh /docker_entrypoint.sh
COPY ./test/platforms/common.sh /common.sh
RUN chmod +x /common.sh
RUN chmod +x /docker_entrypoint.sh

RUN ln -s /pam_ohmage/bin/pam_ohmage.so /usr/lib64/security/pam_ohmage.so

RUN sed -i -E 's/^(\s*)system\(\);/\1unix-stream("\/dev\/log");/' /etc/syslog-ng/syslog-ng.conf

RUN printf ' filter f_pam_ohmage { program("pam_ohmage.*") };\n destination d_pam_ohmage { file("/var/log/pam_ohmage.log"); };\n log { source(s_sys); filter(f_pam_ohmage); destination(d_pam_ohmage); };' >> /etc/syslog-ng/conf.d/pam_ohmage.conf

ENTRYPOINT /docker_entrypoint.sh