package main

import (
	"fmt"
	"github.com/crunchydata/skybridge2/skybridge"
	"os"
)

const TTL = uint64(36000000)
const ETCD = "http://127.0.0.1:4001"

func main() {

	fmt.Println("client starting...")

	hostname := os.Args[1]
	ip := os.Args[2]
	fmt.Println("hostname=" + hostname)
	fmt.Println("ip=" + ip)
	skybridge.AddEntry(hostname, ip, TTL, ETCD)

	fmt.Println("done")
}
