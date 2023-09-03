package ipconfig

// Configuration Interfaces provide methods for operating on map of hostname and active mta
type ConfigServiceIPMap interface {
	Put(string, *HostData)
	Contains(string) bool
	GetValue(string) (*HostData, error)
	GetValues() map[string]*HostData
	RemoveKey(...string)
	Size() int
	Clear()
	IsEmpty() bool
}

// IPListI Interfaces provide methods for operating on list of IP retrieved from mock service
type IPListI interface {
	GetIPList() *IPConfigs
	GetIPValues() []*IPConfigData
	SetIPList([]*IPConfigData)
	IsEmpty() bool
	Clear()
	Size() int
}
