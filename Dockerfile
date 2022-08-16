FROM golang:latest as build

RUN git clone https://github.com/Comradin/yummy.git /yummy

WORKDIR /yummy
RUN cd cmd/server && go mod download && go build -o /yummy/yummy

FROM centos:7
LABEL maintainer Marcus Franke <marcus.franke@gmail.com>

RUN yum -y update && yum install -y createrepo && yum clean all
RUN mkdir -p /usr/share/doc/yummy

COPY .yummy.yml /root
COPY --from=build /yummy/yummy /bin
ADD README.md /usr/share/doc/yummy

EXPOSE 8080

CMD ["yummy"]
