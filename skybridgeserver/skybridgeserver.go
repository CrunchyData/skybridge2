/*
 Copyright 2015 Crunchy Data Solutions, Inc.
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package main

// skybridge is meant to run on any Docker host that
// needs to register DNS entries for any started or stopped container
// command line options (required)
// -d domain name to use (default: crunchy.lab)
// -s skydns client url (e.g. http://192.168.0.106:4100)
// -h docker socket (default: unix://var/run/docker.sock)
// -t TTL value (default: 36000000 )

import (
	//"errors"
	"flag"
	"fmt"
	"github.com/crunchydata/skybridge2/skybridge"
	dockerapi "github.com/fsouza/go-dockerclient"
	"strconv"
	"time"
)

var MAX_TRIES = 3

const delaySeconds = 5
const delay = (delaySeconds * 1000) * time.Millisecond

var DOMAIN string
var ETCD string
var DOCKER_HOST string
var TTL uint64

func init() {
	flag.StringVar(&DOMAIN, "d", "crunchy.lab", "domain name to use when creating DNS entries example crunchy.lab")
	flag.StringVar(&ETCD, "s", "http://127.0.0.1:4001", "URL of etcd client example http://192.168.0.106:4001")
	flag.StringVar(&DOCKER_HOST, "h", "unix:///var/run/docker.sock", "docker socket url")
	flag.Uint64Var(&TTL, "t", 36000000, "dns entries ttl value")
	flag.Parse()
}

func main() {

	var dockerConnected = false
	fmt.Println("DOCKER_HOST=" + DOCKER_HOST)
	fmt.Println("ETCD=" + ETCD)
	fmt.Println("TTL=" + strconv.FormatUint(TTL, 10))
	fmt.Println("DOMAIN=" + DOMAIN)
	var tries = 0
	var docker *dockerapi.Client
	var err error
	for tries = 0; tries < MAX_TRIES; tries++ {
		docker, err = dockerapi.NewClient(DOCKER_HOST)
		err = docker.Ping()
		if err != nil {
			fmt.Println("could not ping docker host")
			fmt.Println("sleeping and will retry in %d sec\n", delaySeconds)
			time.Sleep(delay)
		} else {
			fmt.Println("no err in connecting to docker")
			dockerConnected = true
			break
		}
	}

	if dockerConnected == false {
		fmt.Println("failing, could not connect to docker after retries")
		panic("cant connect to docker")
	}

	events := make(chan *dockerapi.APIEvents)
	assert(docker.AddEventListener(events))
	fmt.Println("skybridge: Listening for Docker events...")
	for msg := range events {
		switch msg.Status {
		//case "start", "create":
		case "start":
			fmt.Println("event: " + msg.Status + " ID=" + msg.ID + " From:" + msg.From)
			skybridge.Action(msg.Status, msg.ID, docker, TTL, ETCD, DOMAIN)
		case "stop":
			fmt.Println("event: " + msg.Status + " ID=" + msg.ID + " From:" + msg.From)
			skybridge.Action(msg.Status, msg.ID, docker, TTL, ETCD, DOMAIN)
		case "destroy":
			fmt.Println("event: " + msg.Status + " ID=" + msg.ID + " From:" + msg.From)
			skybridge.Action(msg.Status, msg.ID, docker, TTL, ETCD, DOMAIN)
		case "die":
			fmt.Println("event: " + msg.Status + " ID=" + msg.ID + " From:" + msg.From)
		default:
			fmt.Println("event: " + msg.Status)
		}
	}

}

func assert(err error) {
	if err != nil {
		fmt.Println("skybridge: ", err)
		panic("can't continue")
	}
}
