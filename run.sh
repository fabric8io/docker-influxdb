#!/bin/bash

sed -i "s|dir = \"/tmp/influxdb/development/state\"|dir = \"${INFLUXDB_DATA_DIR}/state\"|" ${CONFIG_FILE}
sed -i "s|dir = \"/tmp/influxdb/development/db\"|dir = \"${INFLUXDB_DATA_DIR}/db\"|" ${CONFIG_FILE}
sed -i "s|dir  = \"/tmp/influxdb/development/raft\"|dir = \"${INFLUXDB_DATA_DIR}/raft\"|" ${CONFIG_FILE}

mkdir -p ${INFLUXDB_DATA_DIR}
chown -R ${INFLUXDB_USER}:${INFLUXDB_GROUP} ${INFLUXDB_DATA_DIR}

exec sudo -u ${INFLUXDB_USER} -H sh -c "cd /opt/influxdb; exec ./influxd -config ${CONFIG_FILE}"
