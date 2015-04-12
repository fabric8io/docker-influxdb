FROM debian:wheezy

MAINTAINER Jimmi Dyson <jimmidyson@gmail.com>

ENV INFLUXDB_VERSION 0.9.0-rc19
ENV CONFIG_FILE /opt/influxdb/influxdb.conf
ENV INFLUXDB_DATA_DIR /influxdb
ENV INFLUXDB_USER influxdb
ENV INFLUXDB_ADMIN_PORT 8083
ENV INFLUXDB_BROKER_PORT 8086
ENV INFLUXDB_DATA_PORT 8086

RUN apt-get update && apt-get upgrade -y && apt-get install -y wget ca-certificates sudo && rm -rf /var/lib/apt/lists/*

RUN wget -q -O /tmp/influxdb.deb http://get.influxdb.org/influxdb_${INFLUXDB_VERSION}_amd64.deb && \
    dpkg -i /tmp/influxdb.deb && \
    rm -f /tmp/influxdb.deb

RUN cp /etc/opt/influxdb/influxdb.conf /opt/influxdb/influxdb.conf

RUN sed -i 's|^reporting-disabled.*=.*|reporting-disabled = true|' ${CONFIG_FILE} && \
    sed -i 's|file   = "/var/log/influxdb/influxd.log".*||' ${CONFIG_FILE}

EXPOSE 8083 8086 8087

ADD influxdb.conf.tmpl /opt/influxdb/influxdb.conf.tmpl
ADD build/start-influxdb /start-influxdb

CMD ["/start-influxdb"]
