package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/grandcat/zeroconf"
)

const SERVICE = "_workstation._tcp"
const DOMAIN = "local."
const TIMEOUT = 15

type Group struct {
	Hosts []string               `json:"hosts"`
	Vars  map[string]interface{} `json:"vars,omitempty"`
}

func GetInventory() map[string]interface{} {
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		log.Fatal(err)
	}

	inventory := make(map[string]interface{})
	// Channel to receive discovered service entries
	entries := make(chan *zeroconf.ServiceEntry)
	go func(results <-chan *zeroconf.ServiceEntry, inventory map[string]interface{}) {
		hostvars := make(map[string]interface{})

		for entry := range results {
			group := Group{
				Hosts: []string{strings.TrimSuffix(entry.HostName, ".")},
			}
			groupname := strings.ReplaceAll(entry.HostName, ".", "")
			inventory[groupname] = group

			// Building the hostvars for meta
			hvars := make(map[string]interface{})
			hvars["public_ip"] = entry.AddrIPv4[0]
			hostvars[strings.TrimSuffix(entry.HostName, ".")] = hvars
		}

		meta := make(map[string]interface{})
		meta["hostvars"] = hostvars
		inventory["_meta"] = meta
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
