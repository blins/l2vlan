# Driver for connecting containers to external vlan by L2

translated by google.

## Prehistory

It took me to run several containers in the docker swarm and release their traffic directly to a specific vlan.
Everything was ready from the point of view of the classical organization of the network: vlan was created on the router, hung on the gateway
and reach on the switches to the port on the servers. Each container must have its own fixed address.
Without thinking twice, I created macvlan networks on all nodes and launched one container.
All was good. I launched the second container. And then everything was fine. Only here the third container did not want to run ...

Let me dig, what could have gone wrong. The container did not start with an error (approx. Translation): "Gateway (ip) cannot
to be appointed because he is already appointed somewhere. "

Figasse - I think. This address is nowhere to be found, ran through all the interfaces - no, from the word "absolutely." I notice that this
does not work on the site where this container is already running or, suddenly, there was an attempt to launch the container. I'm in absolute misunderstanding
stopped running containers that are tied to these networks and let's remove the networks. The network itself is gone, but the configs
networks do not want to be removed, they say that they are used somewhere. I checked everything, did not find use. I sneaked on their Internet sites:
They will come that this is such a bug and it would be necessary to stop the docker service and delete the local storage ... and then you will be lost
all local configs and you will be happy ... It's good that so far I have only these local configs.

Okay, I started up and began to think. I made a booth on my local computer and reproduced the error.

It is reproduced, by the way, very easily:

```bash
$ docker network create -d macvlan --subnet=192.168.1.0/24 --gateway=192.168.1.1 --ip-range=192.168.1.2/32 \
    --config-only -o parent=eno1.2000 net1conf
$ docker network create -d macvlan --subnet=192.168.1.0/24 --gateway=192.168.1.1 --ip-range=192.168.1.3/32 \
    --config-only -o parent=eno1.2000 net2conf
$ docker network create -d macvlan --scope swarm --config-from net1conf net1 
$ docker network create -d macvlan --scope swarm --config-from net2conf net2
$ docker stack deploy -c docker-compose1.yml stack1
# everything is still good here
$ docker stack deploy -c docker-compose2.yml stack2
# but it was already a mistake
```

The stacks are very simple: they consist of one nginx service, which in one instance is started and connected to an external one.
net1 and net2 respectively.

For example, here’s one file:
```yaml
version: "3.7"

services:
  nginx:
    image: nginx:alpine
    networks:
      - net1
    deploy:
      restart_policy:
        condition: any
      mode: replicated
      replicas: 1

networks:
  net1:
    external: true
```

I downloaded the source code for moby and libnetwork and began to look for what generates such an error. It turned out the driver is IPAM. Such behavior
can be corrected if you write your driver and use it. "Well, OK". No sooner said than done.

I wrote a driver, where the behavior from the standard differs only in that it does not swear at the gateway at all,
he just takes note of it.

And here I am full of anticipation and triumph, I launch everything according to the above example and get a hard break. When creating net2
macvlan driver cannot create interface eno1.2000. Garbage question - use the driver bridge. And here everything is fine
but only this driver for some horseradish hung the gateway address on the used bridge. And there was no way to convince him.

In general, I came to the conclusion that you need to write your network driver. And here he is.

## How it works

The driver creates a bridge with the specified name. Creates a sub-interface with the specified vlan tag to the specified interface.
Combines the created.
And then simply connects the L2 container to this bridge.

If the bridge has already been created, the driver does not check for a sub-interface with the specified vlan tag? because he is supposed to
already created and working.

It is strange why gentlemen programmers from the docker company did not think of this. surely it would be in demand.

## How to use

First you need to install the plugin:
```bash
$ docker plugin install --alias "l2vlan" blins1999/l2vlan
```
Create a network like this:
```bash
$ docker network create -d l2vlan:latest --ipam-driver l2vlan:latest --subnet=192.168.1.0/24 --gateway=192.168.1.1 \
    --ip-range=192.168.1.4/32 -o vlan_id=2000 -o ext_if=eno1 -o bridge_name=vlan2000 --config-only net1conf
$ docker network create -d l2vlan:latest --ipam-driver l2vlan:latest --scope swarm --config-from net1conf net1
```
Needless to say that the first command is executed on EVERY cluster node.

Then you can add static addresses:
```bash
$ docker network create -d l2vlan:latest --ipam-driver l2vlan:latest --subnet=192.168.1.0/24 --gateway=192.168.1.1 \
    --ip-range=192.168.1.6/32 -o vlan_id=2000 -o ext_if=eno1 -o bridge_name=vlan2000 --config-only net2conf
$ docker network create -d l2vlan:latest --ipam-driver l2vlan:latest --scope swarm --config-from net2conf net2
```
And all this will be connected to the same brdige (within one narrow course)

## What to do

  - It is necessary to deal with linux capabilities. I would be happy to help.
 

Well, everything seems :)

