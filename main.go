package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/grandcat/zeroconf"
)

const SERVICE = "_workstation._tcp"
const DOMAIN = "local."
const TIMEOUT = 15

type Group struct {
	Hosts []net.IP               `json:"hosts"`
	Vars  map[string]interface{} `json:"vars"`
}

func GetInventory() map[string]Group {
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		log.Fatal(err)
	}

	inventory := make(map[string]Group)
	// Channel to receive discovered service entries
	entries := make(chan *zeroconf.ServiceEntry)

	go func(results <-chan *zeroconf.ServiceEntry, inventory map[string]Group) {
		for entry := range results {
			vars := make(map[string]interface{})
			group := Group{
				Hosts: entry.AddrIPv4,
				Vars:  vars,
			}
			groupname := strings.ReplaceAll(entry.HostName, ".", "")
			inventory[groupname] = group
		}
	}(entries, inventory)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*TIMEOUT)
	defer cancel()

	err = resolver.Browse(ctx, SERVICE, DOMAIN, entries)
	if err != nil {
		log.Fatalln("Failed to browse:", err.Error())
	}

	<-ctx.Done()

	return inventory
}

func main() {

	isList := flag.Bool("list", false, "List local hosts")
	isHost := flag.Bool("host", false, "get values from an host")
	flag.Parse()

	var inventory interface{}
	if *isHost == false {
		inventory = GetInventory()
	}
	if *isList == false && *isHost == true {
		inventory = make(map[string]interface{})
	}
	b, err := json.Marshal(inventory)
	if err != nil {
		log.Println("error:", err)
	}

	fmt.Println(string(b))
}
