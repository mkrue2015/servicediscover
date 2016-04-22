// Put documentation here
package servicediscover

import (
	"sync"
	"time"

	"github.com/hashicorp/consul/api"
)

// Short description
type Address struct {
	IP   string
	Port int
}

// Short Description
var (
	ConsulServices = make(map[Address]string)
	mutex          sync.Mutex
)

func getServices(consul *api.Client) map[string][]string {
	catalog := consul.Catalog()

	results, _, err := catalog.Services(nil)
	if err != nil {
		return nil
	}

	return results
}

func getServiceInfo(consul *api.Client, entry string, tag string) {
	catalog := consul.Catalog()
	serviceInfo, _, err := catalog.Service(entry, tag, nil)
	if err != nil {
		return
	}
	for _, service := range serviceInfo {
		var vservice Address
		var node Address
		node.IP = service.Address
		node.Port = service.ServicePort
		vservice.IP = service.ServiceAddress
		vservice.Port = service.ServicePort

		ConsulServices[node] = service.ServiceName
		ConsulServices[vservice] = service.ServiceName
	}
}

// Short description
func Services(consulAddress string) {
	config := api.DefaultConfig()
	config.Address = consulAddress
	consul, err := api.NewClient(config)
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

func updateServices(consulAddress string) {
	tick := time.Tick(1 * time.Minute)
	for range tick {
		go func() {
			mutex.Lock()
			defer mutex.Unlock()
			Services(consulAddress)
		}()
	}
}

// Short description
func ServiceName(consulAddress string, ip string, port int) string {
	if len(ConsulServices) == 0 {
		mutex.Lock()
		Services(consulAddress)
		mutex.Unlock()

		go updateServices(consulAddress)
	}

	node := Address{IP: ip, Port: port}
	mutex.Lock()
	defer mutex.Unlock()

	result := ConsulServices[node]

	return result
}
