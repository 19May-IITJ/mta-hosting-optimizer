package ipconfig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIPConfigs(t *testing.T) {
	ip := NewIPConfigList()
	assert.Equal(t, ip.IsEmpty(), true)
	l := make([]*IPConfig, 0)
	l = append(l, &IPConfig{
		Hostname:    "dummy1",
		IPAddresses: "127.0.0.1",
		Status:      true,
	}, &IPConfig{
		Hostname:    "dummy2",
		IPAddresses: "127.0.0.1",
		Status:      false,
	})
	assert.Equal(t, 0, ip.Size())
	t.Run("Positive Test to check SetIPList", func(t *testing.T) {
		ip.SetIPList(l)
		assert.Equal(t, 2, ip.Size())
	})
	ip.Clear()

	assert.Equal(t, true, ip.IsEmpty())
	assert.Equal(t, 0, ip.Size())

}