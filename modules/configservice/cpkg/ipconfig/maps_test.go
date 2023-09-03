package ipconfig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterMap(t *testing.T) {
	mp := NewMap()
	mp.Put("dummymap1", &HostData{
		HostedIP: []string{"127.0.0.1-1", "127.0.0.2-1"},
		ActiveIP: 2,
	})
	mp.Put("dummymap2", &HostData{
		HostedIP: []string{"127.0.0.3-1", "127.0.0.4-0"},
		ActiveIP: 1,
	})
	mp.Put("dummymap3", &HostData{
		HostedIP: []string{"127.0.0.5-1", "127.0.0.6-1", "127.0.0.7-0"},
		ActiveIP: 2,
	})

	expectedMap := map[string]*HostData{
		"dummymap1": &HostData{
			HostedIP: []string{"127.0.0.1-1", "127.0.0.2-1"},
			ActiveIP: 2,
		},
		"dummymap2": &HostData{
			HostedIP: []string{"127.0.0.3-1", "127.0.0.4-0"},
			ActiveIP: 1,
		},
		"dummymap3": &HostData{
			HostedIP: []string{"127.0.0.5-1", "127.0.0.6-1", "127.0.0.7-0"},
			ActiveIP: 2,
		},
	}
	t.Run("Positive Test Case for Maps", func(t *testing.T) {
		assert.Equal(t, expectedMap, mp.GetValues())
		assert.Equal(t, true, mp.Contains("dummymap1"))
		val, err := mp.GetValue("dummymap2")
		assert.NoError(t, err)
		assert.Equal(t, &HostData{
			HostedIP: []string{"127.0.0.3-1", "127.0.0.4-0"},
			ActiveIP: 1,
		}, val)
		assert.Equal(t, 3, mp.Size())
		assert.Equal(t, false, mp.IsEmpty())
		mp.RemoveKey("dummymap3")
		_, err2 := mp.GetValue("dummymap3")
		assert.Error(t, err2)
		mp.Clear()
		assert.Equal(t, 0, mp.Size())
	})
	t.Run("Negative Test Case for Maps", func(t *testing.T) {
		assert.Equal(t, false, mp.Contains("dummymap4"))
		_, err := mp.GetValue("dummymap4")
		assert.Error(t, err)

	})

}
