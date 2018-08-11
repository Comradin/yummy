FROM golang:latest as build

RUN go get -d github.com/Comradin/yummy
RUN go get github.com/golang/dep/cmd/dep

WORKDIR /go/src/github.com/Comradin/yummy

RUN dep ensure
RUN go build


FROM centos:7

LABEL maintainer Marcus Franke <marcus.franke@gmail.com>

USER root

RUN yum -y update
RUN yum install -y createrepo
RUN yum clean all
RUN mkdir -p /usr/share/doc/yummy

COPY .yummy.yml /root
COPY --from=build /go/src/github.com/Comradin/yummy/yummy /bin
ADD README.md /usr/share/doc/yummy

EXPOSE 8080

CMD ["yummy", "serve"]
