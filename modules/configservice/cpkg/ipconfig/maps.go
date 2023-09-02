package ipconfig

import (
	"github.com/pkg/errors"
)

var _ ConfigServiceIPMap = (*RegisterMap)(nil)

type RegisterMap struct {
	handlermap map[string]*HostData
}

type IPState struct {
	IP    string
	State bool
}
type HostData struct {
	HostedIP []string
	ActiveIP int
}

// New returns the new map of RegisterMap type
func NewMap() *RegisterMap {
	r := RegisterMap{handlermap: make(map[string]*HostData)}
	return &r
}

// Put Adds the key value pair into the map
func (s *RegisterMap) Put(key string, value *HostData) {
	s.handlermap[key] = value
}

// Contains check whehter the key is present in map
func (s *RegisterMap) Contains(key string) bool {
	_, exists := s.handlermap[key]
	return exists
}

// GetValues returns value associated with the key
func (s *RegisterMap) GetValue(key string) (values *HostData, err error) {
	var ok bool
	values, ok = s.handlermap[key]
	if !ok {
		return values, errors.Errorf("no value for key %v", key)
	}
	return values, nil
}

// RemoveKey remove value associated with the key
func (s *RegisterMap) RemoveKey(keys ...string) {
	for _, key := range keys {
		delete(s.handlermap, key)
	}
}

// GetValues returns the underlying map
func (s *RegisterMap) GetValues() map[string]*HostData {
	return s.handlermap
}

// Size returns the no. of entries in map
func (s *RegisterMap) Size() int {
	return len(s.handlermap)
}

// Flushes all the key in map
func (s *RegisterMap) Clear() {
	s.handlermap = make(map[string]*HostData)
}

// Check whether the Map is empty or not
func (s *RegisterMap) IsEmpty() bool {
	return s.Size() == 0
}
