FROM centos:7

LABEL maintainer Marcus Franke <marcus.franke@gmail.com>

USER root

RUN yum -y update
RUN yum install -y createrepo
RUN yum clean all

COPY .yummy.yml /root
ADD yummy /bin

EXPOSE 8080

CMD ["yummy", "serve"]
