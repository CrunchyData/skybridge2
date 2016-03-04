

skybridge
===========

requires skydns 2.0.1d
requires etcd 2.0.0

skybridge is a small program that listens to Docker events locally
and adds/removes/updates DNS information on a skydns service
either locally or remotely. RHEL7/Centos7 64 bit systems are
currently supported.

skybridge is distributed as a Docker image found on github
at the following location:

https://registry.hub.docker.com/u/crunchydata/skybridge

There is a script to run the skybridge docker container
found here:

https://github.com/CrunchyData/skybridge/blob/master/bin/run-skybridge.sh


environment prerequisites
=========================

First, you will need a stable hostname and IP address, so
set your hostname (e.g. dev.crunchy.lab) using the nmtui utility
or similar.

Second, make sure to add your hostname to your /etc/hosts file!

Next, you are setting up a DNS server by installing skybridge!  So,
you will need a basic understanding of Linux networking, but
here are the basics you will need:

Make sure /etc/resolv.conf specifies your server's IP address as the primary DNS nameserver as follows, using 192.168.0.106 as an example of 
your server's IP address and 192.168.0.1 as your existing DNS nameserver:

~~~~~~~~~~~~~~~~~~
nameserver 192.168.0.106
nameserver 192.168.0.1
~~~~~~~~~~~~~~~~~~

You will want to adjust your network settings to ensure that these
DNS values don't get blown away after you reboot your system! See
PEERDNS and NetworkManager for details on how to set your
ethernet adapter settings. A static IP address is essential
for a production installation, however if you just want to test
skybridge, and your DHCP address doesn't change very often, you 
can specify your DHCP assigned address.

Docker Container Installation
=============================
The easy way to run skybridge is to run the Docker container version.

A container has been built that includes etcd, skydns, and skybridge
ready for use.

The container is located in DockerHub at crunchydata/skybridge:latest

To execute, run the run-skybridge.sh script found here:
https://github.com/CrunchyData/skybridge/blob/master/bin/run-skybridge.sh

Edit the script, adding your own IP address of your host, the domain
name of your choice.

Then run the script:

sudo ./run-skybridge.sh

This script will pull down the skybridge docker image, and execute
it.

Host Installation
=====================
Users can also run skybridge on their host, outside of a container, by
downloading the skybridge installation archive from
the following location:

~~~~~~~~~~~~~~~~~~~~~~~~~
wget https://s3.amazonaws.com/crunchydata/cpm/skybridge.1.0.4-linux-amd64.tar.gz
~~~~~~~~~~~~~~~~~~~~~~~~~

They will un-tar the file and run the install.sh script located
inside the archive.  The install.sh script will prompt them
through the install. The installation script requires sudo
privileges.

~~~~~~~~~~~~~~~~~~~~~~~~~
tar xvzf skybridge.1.0.0-linux-amd64.tar.gz
./install.sh
~~~~~~~~~~~~~~~~~~~~~~~~~

The install.sh script will prompt the user for the IP address to
use for running the etcd/skydns/skybridge services, as well
as the domain name to use.

systemd unit files are copied to the user's system (/usr/lib/systemd/system), enabled, and started.

All installed files are copied to the following directory:
~~~~~~~~~~~~~~~~~~~~~~~~~
/var/cpm/bin
/var/cpm/config
~~~~~~~~~~~~~~~~~~~~~~~~~


Building from Source
==========================

Here are steps to build skybridge2 from source:

~~~~
mkdir -p sky/src sky/bin sky/pkg
export GOPATH=$HOME/sky;export GOBIN=$GOPATH/bin;export PATH=$PATH:$GOBIN
cd sky
go get github.com/tools/godep
go get github.com/crunchydata/skybridge2
cd src/github.com/crunchydata/skybridge2
godep restore
cd /tmp
tar -xzvf $GOPATH/src/github.com/crunchydata/skybridge2/archives/etcd-2.0.0.tar.gz ./etcd-v2.0.0-linux-amd64/etcdctl
cp ./etcd-v2.0.0-linux-amd64/etcdctl $GOBIN
tar -xzvf $GOPATH/src/github.com/crunchydata/skybridge2/archives/etcd-2.0.0.tar.gz ./etcd-v2.0.0-linux-amd64/etcd
cp ./etcd-v2.0.0-linux-amd64/etcd $GOBIN
tar -xzvf $GOPATH/src/github.com/crunchydata/skybridge2/archives/skydns-2.0.1d.tar.gz ./bin/skydns
cp ./bin/skydns $GOBIN
cd $GOPATH/src/github.com/crunchydata/skybridge2
make build
make image
~~~~


etcd
===========

etcd is included in the user installation archive.

etcd is used to store DNS information for skydns, in a typical
deployment, you would install etcd as a first step

Install instructions for etcd are found at https://github.com/coreos/etcd/releases

Currently we are using etcd 2.0.0.

By default, etcd binds to localhost.

We specify a location for etcd to store it's data, for
example:

~~~~~~~~~~~~~~~~~~
-data-dir /var/cpm/data/etcd
~~~~~~~~~~~~~~~~~~


Starting etcd
-----------------

A systemd unit file is provided, etcd.service, for automatically
starting etcd, see config/etcd.service for an example.  This file
is installed when a user performs a skybridge user install using
the binary tar archive.

skydns
=================
skydns is a DNS server that we are using specifically to
discover Docker container instances.

Skydns uses etcd to store it's data and is therefore a dependency.

We are currently using skydns 2.0.1d which can be found at:
~~~~~~~~~~~~~~~~~~
git clone git@github.com:skynetservices/skydns.git
git checkout tags/2.0.1d
~~~~~~~~~~~~~~~~~~

After building skydns, you will pass a flag to it which identifies
the backend etcd system it will use, for example:
~~~~~~~~~~~~~~~~~~
export ETCD_MACHINES='http://192.168.0.106:4001'
~~~~~~~~~~~~~~~~~~
or
~~~~~~~~~~~~~~~~~~
-machines=127.0.0.1:4001
~~~~~~~~~~~~~~~~~~

Starting skydns
-------------------

Remember that DNS is a privileged port and requires you start
skydns as root if you want to use the default port of DNS (53)

A systemd unit file is provided to start skydns, it is found in
the config/skydns.service file.

After starting skydns, you can test it using curl.
Example of adding a host and IP address to skydns using curl:
~~~~~~~~~~~~~~~~~~
curl -XPUT http://127.0.0.1:4001/v2/keys/skydns/lab/crunchy/foo \
-d value='{"host":"192.168.0.107", "port":8080}'
dig foo.crunchy.lab
~~~~~~~~~~~~~~~~~~

skybridge
===================

skybridge is meant to be installed on any Docker host.  skybridge
will listen to the local Docker service, once a start or stop Docker
event is received, skybridge opens a connection to the 
skydns server and makes a skydns REST API call to create or update
a DNS record.

So, if a Docker container called pgdb1 is started, skybridge
will send to skydns the container's IP Address, domain name, and
container name.

skybridge configuration flags
----------------------------

DOMAIN
------
this is the DNS domain name to use when registering new containers
or containers that have been removed, the default is 'crunchy.lab' if
not specified, for example:
~~~~~~~~~~~~~~~
-d crunchy.lab
~~~~~~~~~~~~~~~

SKYDNS
------
this is the URL of the etcd server to be used for registering
DNS information, the default is http://127.0.0.1:4001, if not
specified:
~~~~~~~~~~~~~~~
-s http://192.168.0.106:4001
~~~~~~~~~~~~~~~

DOCKER
------
this is the URL of the Docker server socket that will be listened to
by skybridge, the value of the env variable SWARM_MANAGER_URL is used:
~~~~~~~~~~~~~~~
-h tcp://192.168.0.107:8000
~~~~~~~~~~~~~~~

TTL
------
this is the TTLS value to use when registering DNS values, the default
is 360 if not specified by this flag:
~~~~~~~~~~~~~~~
-t 400
~~~~~~~~~~~~~~~

Starting skybridge
-------------------------

A systemd unit file is provided to start skybridge, it is found in
the config/skybridge.service file.

Example startup of skybridge that will connect to an etcd service
running at 192.168.0.106:4001 and use a domain name of 'crunchy.lab'
for creating new entries in DNS:
~~~~~~~~~~~~~~~~~~
./skybridge -s http://192.168.0.106:4001  -d "crunchy.lab."
~~~~~~~~~~~~~~~~~~

New entries take the form of containerName.domainname


Testing the User Installation
=============================

After the skybridge installation has been performed, you
can create a Docker container and then query the DNS name as
follows:
~~~~~~~~~~~~~~~~~~~~~~~~~
docker run -d --name=tester busybox /bin/sh
~~~~~~~~~~~~~~~~~~~~~~~~~

In another terminal window:
~~~~~~~~~~~~~~~~~~~~~~~~~
dig tester.crunchy.lab
~~~~~~~~~~~~~~~~~~~~~~~~~

If all is working as normal, the tester.crunchy.lab should resolv
to your container's IP address.  You can also do a reverse
DNS lookup using the IP address as follows:

~~~~~~~~~~~~~~~~~~~~~~~~~
dig 172.17.0.XXX
~~~~~~~~~~~~~~~~~~~~~~~~~


Manual Adding a DNS Entry
=============================

Within the container is the /var/cpm/bin/skybridgeclient binary.
With this command you can create a DNS entry manually which is
useful for testing, first parameter is the hostname without the
domain name, and the second parameter is the ip address of that
host:

~~~~
skybridgeclient <somehostname> <some ip address>
~~~~
