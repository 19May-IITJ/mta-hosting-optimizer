package ipconfig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHostnamesWithMaxIPs(t *testing.T) {
	mp := NewMap()
	mp.Put("dummy1", 2)
	mp.Put("dummy2", 1)
	mp.Put("dummy3", 0)
	inefficientHostname := GetHostnamesWithMaxIPs(1, mp)
	t.Run("Positive Test for GetHostnamesWithMaxIPs", func(t *testing.T) {
		expectPositivehostname := []string{"dummy2", "dummy3"}
		assert.Equal(t, expectPositivehostname, inefficientHostname)
	})

	t.Run("Negative Test for GetHostnamesWithMaxIPs", func(t *testing.T) {
		expectNegativeTestResulthostname := []string{"dummy1", "dummy3"}
		assert.NotEqual(t, expectNegativeTestResulthostname, inefficientHostname)
	})

}
