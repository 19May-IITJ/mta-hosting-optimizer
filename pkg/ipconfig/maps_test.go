package ipconfig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterMap(t *testing.T) {
	mp := NewMap()
	mp.Put("dummymap1", 2)
	mp.Put("dummymap2", 1)
	mp.Put("dummymap3", 0)
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
		assert.Equal(t, 0, val)
		mp.Clear()
		assert.Equal(t, 0, mp.Size())
	})
	t.Run("Negative Test Case for Maps", func(t *testing.T) {
		assert.Equal(t, false, mp.Contains("dummymap4"))
		val, err := mp.GetValue("dummymap4")
		assert.Error(t, err)
		assert.Equal(t, 0, val)
	})

}
