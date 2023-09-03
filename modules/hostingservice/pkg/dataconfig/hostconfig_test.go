package dataconfig

import (
	"mta2/modules/utility"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterMap(t *testing.T) {
	mp := NewHostMap()

	hd := make([]*utility.Message, 0)
	hd = append(hd, &utility.Message{
		Hostname: "dummymap1",
		Active:   2,
	}, &utility.Message{
		Hostname: "dummymap2",
		Active:   1,
	},
		&utility.Message{
			Hostname: "dummymap3",
			Active:   0,
		})

	mp.Put(hd...)
	t.Run("Positive Test Case for Maps", func(t *testing.T) {
		assert.Equal(t, true, mp.Contains("dummymap1"))
		val, err := mp.GetValue("dummymap2")
		assert.NoError(t, err)
		assert.Equal(t, 1, val)
		assert.Equal(t, 3, mp.Size())
		assert.Equal(t, false, mp.IsEmpty())
		mp.RemoveKey("dummymap3")
		val, err = mp.GetValue("dummymap3")
		assert.Error(t, err)
		assert.Equal(t, -1, val)
		mp.Clear()
		assert.Equal(t, 0, mp.Size())
	})
	t.Run("Negative Test Case for Maps", func(t *testing.T) {
		assert.Equal(t, false, mp.Contains("dummymap4"))
		val, err := mp.GetValue("dummymap4")
		assert.Error(t, err)
		assert.Equal(t, -1, val)
	})

}
