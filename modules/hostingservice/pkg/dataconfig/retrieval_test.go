package dataconfig

import (
	"mta2/modules/utility"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHostnamesWithMaxIPs(t *testing.T) {
	mp := NewHostMap()

	hd := make([]*utility.Message, 0)
	hd = append(hd, &utility.Message{
		Hostname: "dummy1",
		Active:   2,
	}, &utility.Message{
		Hostname: "dummy2",
		Active:   1,
	},
		&utility.Message{
			Hostname: "dummy3",
			Active:   0,
		})

	mp.Put(hd...)

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
