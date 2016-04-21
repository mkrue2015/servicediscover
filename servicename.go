// Put documentation here
package servicediscover

import (
	consulapi "github.com/hashicorp/consul/api"
)

type Address struct {
	Ip   string
	Port int
}

var ConsulServices map[Address]string = make(map[Address]string)

// Short description
func getServices(consul *consulapi.Client) map[string][]string {
	catalog := consul.Catalog()

	results, _, err := catalog.Services(nil)
	if err != nil {
		return nil
	}

	return results
}

// Short description
func getServiceInfo(consul *consulapi.Client, entry string, tag string) {
	catalog := consul.Catalog()
	serviceInfo, _, err := catalog.Service(entry, tag, nil)
	if err != nil {
		return
	}
	for _, service := range serviceInfo {
		var node Address
		node.Ip = service.ServiceAddress
		node.Port = service.ServicePort

		ConsulServices[node] = service.ServiceName
	}
}

// Short description
func Services(consulAddress string) {
	config := consulapi.DefaultConfig()
	config.Address = consulAddress
	consul, err := consulapi.NewClient(config)
	if err != nil {
		return
	}

	var catalogEntries map[string][]string

	catalogEntries = getServices(consul)

	for catalogEntry, catalogTags := range catalogEntries {
		for _, tag := range catalogTags {
			getServiceInfo(consul, catalogEntry, tag)
		}
	}
}

// Short description
func ServiceName(consulAddress string, ip string, port int) string {
	Services(consulAddress)

	node := Address{Ip: ip, Port: port}
	return ConsulServices[node]
}
