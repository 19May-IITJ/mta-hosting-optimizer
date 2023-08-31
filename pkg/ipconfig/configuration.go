package ipconfig

// Configuration Interfaces provide methods for operating on map of hostname and active mta
type Configuration interface {
	Put(string, int)
	Contains(string) bool
	GetValue(string) (int, error)
	GetValues() map[string]int
	RemoveKey(...string)
	Size() int
	Clear()
	IsEmpty() bool
}

// IPListI Interfaces provide methods for operating on list of IP retrieved from mock service
type IPListI interface {
	GetIPList() *IPConfigs
	GetIPValues() []*IPConfig
	SetIPList([]*IPConfig)
	IsEmpty() bool
	Clear()
	Size() int
}
