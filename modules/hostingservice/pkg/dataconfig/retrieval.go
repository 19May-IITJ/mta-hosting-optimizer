package dataconfig

import "sync"

var DataMutex sync.RWMutex

// GetHostnamesWithMaxIPs checks which hostname has active MTA less than threshold value and return corresponding hostnames
func GetHostnamesWithMaxIPs(maxIPs int, hostmap HostingServiceHostMap) []string {
	DataMutex.Lock()
	defer DataMutex.Unlock()

	result := make([]string, 0)
	for hostname, activeMTAs := range hostmap.GetValues() {
		if activeMTAs <= maxIPs {
			result = append(result, hostname)
		}
	}
	return result
}
