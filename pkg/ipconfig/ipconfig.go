package ipconfig

import (
	"sync"
)

type IPConfig struct {
	Hostname    string `json:"hostname"`
	IPAddresses string `json:"ipAddresses"`
	Status      bool   `json:"status"`
}

var dataMutex sync.RWMutex

func GetHostnamesWithMaxIPs(maxIPs int, iplist Configuration) []string {
	dataMutex.RLock()
	defer dataMutex.RUnlock()

	result := make([]string, 0)
	for hostname, activeMTAs := range iplist.GetValues() {
		if activeMTAs <= maxIPs {
			result = append(result, hostname)
		}
	}
	return result
}
