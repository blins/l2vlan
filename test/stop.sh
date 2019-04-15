#!/usr/bin/env bash

set -x

docker stack rm stack1 stack2

sleep 5

docker network rm net1 net2

sleep 5

docker network rm net1conf net2conf
