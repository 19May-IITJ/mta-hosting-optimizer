package dataconfig

import "mta2/modules/utility"

type HostingServiceHostMap interface {
	Put(...*utility.Message)
	Contains(string) bool
	GetValue(string) (int, error)
	GetValues() map[string]int
	RemoveKey(...string)
	Size() int
	Clear()
	IsEmpty() bool
}
