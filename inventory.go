package main

import (
	"net"
)

type Group struct {
	Hosts []net.IP               `json:"hosts"`
	Vars  map[string]interface{} `json:"vars,omitempty"`
}

type HostVars struct {
	Key   string
	Value interface{}
}

type Inventory struct {
	Grp map[string]Group
	Hv  map[string][]HostVars
}

func NewInventory() *Inventory {
	inv := new(Inventory)
	inv.Grp = make(map[string]Group)
	inv.Hv = make(map[string][]HostVars)

	return inv
}

func (inv Inventory) Build() map[string]interface{} {
	result := make(map[string]interface{})

	// Building meta data
	hostvars := make(map[string]interface{})
	for h, hv := range inv.Hv {
		tmp := make(map[string]interface{})
		for _, v := range hv {
			tmp[v.Key] = v.Value
		}
		hostvars[h] = tmp
	}

	result["_meta"] = map[string]interface{}{
		"hostvars": hostvars,
	}

	// Building group data
	for host, group := range inv.Grp {
		result[host] = group
	}

	return result
}

func (inv Inventory) AddGroup(host string, group Group) Inventory {
	inv.Grp[host] = group
	return inv
}

func (inv Inventory) AddMetaHostvars(host string, key string, value interface{}) Inventory {
	inv.Hv[host] = append(inv.Hv[host], HostVars{Key: key, Value: value})
	return inv
}
