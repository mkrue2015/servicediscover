// Put documentation here
package servicediscover

import (
	"time"
  "sync"

	consulapi "github.com/hashicorp/consul/api"
)

// Short description
type Address struct {
	Ip   string
	Port int
}

// Short Description
var ConsulServices map[Address]string = make(map[Address]string)
var ConsulQueryCount int = 0
var mutex sync.Mutex

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

func updateServices(consulAddress string) {

	tick := time.NewTicker(1 * time.Minute)
  for now := <-tick {
  	go func() {
      mutex.Lock()
      defer mutex.Unlock()
      Services(consulAddress)}()
	}
}

// Short description
func ServiceName(consulAddress string, ip string, port int) string {
	if len(ConsulServices) == 0 {
		go updateServices(consulAddress)
	}

	node := Address{Ip: ip, Port: port}
	ConsulQueryCount++
  mutex.Lock()
  defer mutex.Unlock()

	result :=  ConsulServices[node]

  return result
}
