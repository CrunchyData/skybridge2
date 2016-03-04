package skybridge

import (
	"encoding/json"
	"fmt"
	"github.com/coreos/go-etcd/etcd"
	"github.com/skynetservices/skydns/msg"
	"strings"
)

//global TTL
//global skydns url

//adds a service entry and a PTR entry
func AddEntry(hostname string, ip string, TTL uint64, ETCD string) {

	fmt.Println("addEntry called")

	var services = []*msg.Service{
		{Host: ip, Key: hostname + "."},
		{Host: hostname, Key: reverseIP(ip)},
	}

	client := etcd.NewClient([]string{ETCD})

	//delete any existing entries with this name or ip address
	DeleteEntry(hostname, ip, ETCD)

	//add a service

	fmt.Println("creating A record...")
	serv := services[0]
	b, err := json.Marshal(serv)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	path, _ := msg.PathWithWildcard(serv.Key)

	_, err = client.Create(path, string(b), TTL)
	if err != nil {
		// TODO(miek): allow for existing keys...
		fmt.Println(err.Error())
	}

	//add a PTR
	fmt.Println("creating PTR record...")
	serv = services[1]
	b, err = json.Marshal(serv)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	path, _ = msg.PathWithWildcard(serv.Key)

	_, err = client.Create(path, string(b), TTL)
	if err != nil {
		// TODO(miek): allow for existing keys...
		fmt.Println(err.Error())
	}

	fmt.Println("AddEntry completed")

}

//delete both the service entry and the PTR entry
func DeleteEntry(hostname string, ip string, ETCD string) {
	fmt.Println("DeleteEntry called...")
	var services = []*msg.Service{
		{Host: ip, Key: hostname + "."},
		{Host: hostname, Key: reverseIP(ip)},
	}

	client := etcd.NewClient([]string{ETCD})
	//delete a service

	serv := services[0]
	path, _ := msg.PathWithWildcard(serv.Key)

	_, err := client.Delete(path, false)
	if err != nil {
		fmt.Println(err.Error())
	}

	//delete a PTR

	serv = services[1]
	path, _ = msg.PathWithWildcard(serv.Key)

	_, err = client.Delete(path, false)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("DeleteEntry completed...")

}

//return the reverse ip
func reverseIP(ip string) string {
	//"1.0.0.10.in-addr.arpa."},
	//assume ip has 4 numbers 1.2.3.4
	arr := strings.Split(ip, ".")
	return arr[3] + "." + arr[2] + "." + arr[1] + "." + arr[0] + ".in-addr.arpa"
}
