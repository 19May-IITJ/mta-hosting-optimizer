package ipconfig

import (
	"sync"
)

type IPConfigData struct {
	Hostname    string `json:"hostname"`
	IPAddresses string `json:"ipAddresses"`
	Status      bool   `json:"status"`
}

var dataMutex sync.RWMutex

// GetHostnamesWithMaxIPs checks which hostname has active MTA less than threshold value and return corresponding hostnames
func GetHostnamesWithMaxIPs(maxIPs int, iplist ConfigServiceIPMap) []string {
	dataMutex.RLock()
	defer dataMutex.RUnlock()

	result := make([]string, 0)
	// for hostname, activeMTAs := range iplist.GetValues() {
	// 	if activeMTAs <= maxIPs {
	// 		result = append(result, hostname)
	// 	}
	// }
	return result
}
