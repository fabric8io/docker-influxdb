FROM centos:7

MAINTAINER Jimmi Dyson <jimmidyson@gmail.com>

ENTRYPOINT ["/start-influxdb"]
EXPOSE 8083 8086 8087
VOLUME /data

ENV INFLUXDB_VERSION 0.9.0-1
ENV CONFIG_FILE /opt/influxdb/influxdb.conf

RUN rpm -ivh https://s3.amazonaws.com/influxdb/influxdb-0.9.0-1.x86_64.rpm

ADD build/start-influxdb /start-influxdb
ADD influxdb.conf.tmpl /opt/influxdb/influxdb.conf.tmpl

