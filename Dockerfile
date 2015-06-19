FROM centos:7
MAINTAINER Jimmi Dyson <jimmidyson@gmail.com>

# Install InfluxDB
ENV INFLUXDB_VERSION 0.8.8
RUN rpm -ivh https://s3.amazonaws.com/influxdb/influxdb-0.8.8-1.x86_64.rpm

ADD config.toml /config/config.toml
ADD run.sh /run.sh
RUN chmod 777 /config

ENV PRE_CREATE_DB **None**
ENV SSL_SUPPORT **False**
ENV SSL_CERT **None**

# Admin server
EXPOSE 8083

# HTTP API
EXPOSE 8086

# HTTPS API
EXPOSE 8084

# Raft port (for clustering, don't expose publicly!)
#EXPOSE 8090

# Protobuf port (for clustering, don't expose publicly!)
#EXPOSE 8099

VOLUME ["/data"]

ENTRYPOINT ["/run.sh"]
