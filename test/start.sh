#!/usr/bin/env bash

set -x

PLUGIN=blins1999/l2vlan:latest
PREFIX=192.168.1.
ETH=eno1
VLAN=2000

docker plugin enable ${PLUGIN}

sleep 5

docker network create -d ${PLUGIN} --ipam-driver=${PLUGIN}  --subnet=${PREFIX}0/24 --gateway=${PREFIX}1 --ip-range=${PREFIX}2/32 \
    -o vlan_id=${VLAN} -o ext_if=${ETH} -o bridge_name=vlan${VLAN} --config-only net1conf

docker network create -d ${PLUGIN} --ipam-driver=${PLUGIN}  --subnet=${PREFIX}0/24 --gateway=${PREFIX}1 --ip-range=${PREFIX}3/32 \
    -o vlan_id=${VLAN} -o ext_if=${ETH} -o bridge_name=vlan${VLAN} --config-only net2conf

sleep 3

docker network create -d ${PLUGIN} --scope swarm --config-from net1conf net1
docker network create -d ${PLUGIN} --scope swarm --config-from net2conf net2

sleep 3

docker stack deploy -c docker-compose1.yml stack1
docker stack deploy -c docker-compose2.yml stack2