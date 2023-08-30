package ipconfig

import (
	"github.com/pkg/errors"
)

var _ Configuration = (*RegisterMap)(nil)

type RegisterMap struct {
	handlermap map[string]int
}

// New returns the new map of RegisterMap type
func NewMap() *RegisterMap {
	r := RegisterMap{handlermap: make(map[string]int)}
	return &r
}

// Put Adds the key value pair into the map
func (s *RegisterMap) Put(key string, value int) {
	s.handlermap[key] = value
}

// Contains check whehter the key is present in map
func (s *RegisterMap) Contains(key string) bool {
	_, exists := s.handlermap[key]
	return exists
}

// GetValues returns value associated with the key
func (s *RegisterMap) GetValue(key string) (values int, err error) {
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
func (s *RegisterMap) GetValues() map[string]int {
	return s.handlermap
}
func (s *RegisterMap) Size() int {
	return len(s.handlermap)
}
func (s *RegisterMap) Clear() {
	s.handlermap = make(map[string]int)
}
func (s *RegisterMap) IsEmpty() bool {
	return s.Size() == 0
}

// func (s *RegisterMap) GetRandomValue() (value int) {
// 	if !s.IsEmpty() {
// 		random := rand.Intn(s.Size())
// 		for _, val := range s.handlermap {
// 			random--
// 			if random == 0 {
// 				return val
// 			}
// 		}
// 		return

// 	}
// 	return
// }

// func (s *RegisterMap) Copy() RegisterMap {
// 	new := NewMap()
// 	for key, val := range s.GetValues() {
// 		new.Put(key, val)
// 	}
// 	return *new
// }

// func (s *RegisterMap)
