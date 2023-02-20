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

func GetInventory() *Inventory {

	inventory := NewInventory()

	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		log.Fatal(err)
	}

	// Channel to receive discovered service entries
	entries := make(chan *zeroconf.ServiceEntry)
	go func(results <-chan *zeroconf.ServiceEntry, inventory *Inventory) {
		for entry := range results {
			cleanhostname := strings.TrimSuffix(entry.HostName, ".")
			group := Group{
				Hosts: []string{cleanhostname},
			}
			inventory.AddGroup(strings.ReplaceAll(cleanhostname, ".", "_"), group)
			inventory.AddMetaHostvars(cleanhostname, "public_ip", entry.AddrIPv4[0])
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
		inventory = GetInventory().Build()
	} else if *isList == false && *isHost == true {
		inventory = make(map[string]interface{})
	}

	b, err := json.Marshal(inventory)
	if err != nil {
		log.Println("error:", err)
	}

	fmt.Println(string(b))
}
