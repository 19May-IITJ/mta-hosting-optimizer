package dataconfig

import (
	"errors"
	"mta2/modules/utility"
)

type HostMap struct {
	handlermap map[string]int
}

func NewHostMap() *HostMap {
	return &HostMap{
		handlermap: make(map[string]int),
	}
}

var _ HostingServiceHostMap = (*HostMap)(nil)

func (h *HostMap) Put(hd ...*utility.Message) {
	for _, data := range hd {
		h.handlermap[data.Hostname] = data.Active
	}
}

func (h *HostMap) Contains(name string) bool {
	_, found := h.handlermap[name]
	return found
}

func (h *HostMap) GetValue(name string) (int, error) {
	if h.Contains(name) {
		return h.handlermap[name], nil
	}
	return -1, errors.New("key not present")
}

func (h *HostMap) GetValues() map[string]int {
	return h.handlermap
}

func (h *HostMap) RemoveKey(names ...string) {

	for _, name := range names {
		if h.Contains(name) {
			delete(h.handlermap, name)
		}
	}
}

func (h *HostMap) Size() int {
	return len(h.handlermap)
}

func (h *HostMap) Clear() {
	h.handlermap = make(map[string]int)
}

func (h *HostMap) IsEmpty() bool {
	return h.Size() == 0
}
