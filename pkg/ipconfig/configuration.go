package ipconfig

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

type IPListI interface {
	GetIPList() *IPConfigs
	SetIPList([]*IPConfig)
	IsEmpty() bool
	Clear()
	Size() int
}
