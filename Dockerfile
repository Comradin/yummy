FROM centos:7

LABEL maintainer Marcus Franke <marcus.franke@gmail.com>

USER root

RUN yum -y update
RUN yum install -y createrepo
RUN yum clean all
RUN mkdir -p /usr/share/doc/yummy

COPY .yummy.yml /root
ADD yummy /bin
ADD README.md /usr/share/doc/yummy

EXPOSE 8080

CMD ["yummy", "serve"]
