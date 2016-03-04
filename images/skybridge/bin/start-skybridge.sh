#!/bin/bash -x

# Copyright 2016 Crunchy Data Solutions, Inc.
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

#
# start up skybridge, skydns, and etcd
#
# the service looks for the following env vars to be set on
# startup
#
# $DNS_DOMAIN domain name we are serving with DNS
# $DNS_NAMESERVER secondary DNS nameservers to add
#

export SWARM_MANAGER_URL=$SWARM_MANAGER_URL
export DOCKER_HOST=$DOCKER_HOST
export DNS_DOMAIN=$DNS_DOMAIN
export DNS_NAMESERVER=$DNS_NAMESERVER

echo DNS_DOMAIN=$DNS_DOMAIN
echo DNS_NAMESERVER=$DNS_NAMESERVER

export PATH=/var/cpm/bin:$PATH

#
# start etcd
#

etcd --data-dir /etcddata &

sleep 3

#
# start skydns
#

skydns -addr=0.0.0.0:53 -machines=127.0.0.1:4001 -domain=$DNS_DOMAIN. -nameservers=$DNS_NAMESERVER:53 &

sleep 3

#
# start skybridge
#

ls -l /tmp/docker*

#skybridge -d $DNS_DOMAIN -h $SWARM_MANAGER_URL -s http://127.0.0.1:4001 
skybridgeserver -d $DNS_DOMAIN -h unix:///tmp/docker.sock -s http://127.0.0.1:4001 
#skybridge -d $DNS_DOMAIN  -s http://127.0.0.1:4001 

