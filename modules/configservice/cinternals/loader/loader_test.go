package loader

import (
	"context"
	"encoding/json"
	"mta2/mock/mocking"
	"mta2/modules/configservice/cinternals/constants"
	"mta2/modules/configservice/cpkg/ipconfig"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfigIPConfiguration(t *testing.T) {
	// Create mock ipconfig.Configuration and ipconfig.IPList objects
	mockConfig := ipconfig.NewMap()
	mockIPList := ipconfig.NewIPConfigList()

	t.Run("Test with valid JSON file path", func(t *testing.T) {
		os.Setenv(constants.DBPATH, "/tmp/code/mock/test_data/ipconfig_test.json")
		defer os.Unsetenv(constants.DBPATH)
		err := LoadConfigIPConfiguration(mockConfig, mockIPList)
		assert.NoError(t, err)
		mockConfig_expected := ipconfig.NewMap()
		mockIPList_expected := ipconfig.NewIPConfigList()
		mockConfig_expected.Put("dummy_1", &ipconfig.HostData{
			HostedIP: []string{"127.0.0.1-1"},
			ActiveIP: 1,
		})
		mockConfig_expected.Put("dummy_2", &ipconfig.HostData{
			HostedIP: []string{"127.0.0.2-0"},
			ActiveIP: 0,
		})
		mockConfig_expected.Put("dummy_3", &ipconfig.HostData{
			HostedIP: []string{"127.0.0.3-0", "127.0.0.4-1", "127.0.0.5-1"},
			ActiveIP: 2,
		})

		l := make([]*ipconfig.IPConfigData, 0)
		l = append(l,
			&ipconfig.IPConfigData{
				Hostname:    "dummy_1",
				IPAddresses: "127.0.0.1",
				Status:      true,
			},
			&ipconfig.IPConfigData{
				Hostname:    "dummy_2",
				IPAddresses: "127.0.0.2",
				Status:      false,
			},
			&ipconfig.IPConfigData{
				Hostname:    "dummy_3",
				IPAddresses: "127.0.0.3",
				Status:      false,
			},
			&ipconfig.IPConfigData{
				Hostname:    "dummy_3",
				IPAddresses: "127.0.0.4",
				Status:      true,
			},
			&ipconfig.IPConfigData{
				Hostname:    "dummy_3",
				IPAddresses: "127.0.0.5",
				Status:      true,
			},
		)
		mockIPList_expected.SetIPList(l)

		sort.Slice(mockIPList_expected.GetIPValues(), func(i, j int) bool {
			return mockIPList_expected.GetIPValues()[i].Hostname < mockIPList_expected.GetIPValues()[j].Hostname
		})
		sort.Slice(mockIPList.GetIPValues(), func(i, j int) bool {
			return mockIPList.GetIPValues()[i].Hostname < mockIPList.GetIPValues()[j].Hostname
		})
		assert.Equal(t, mockConfig_expected, mockConfig)
		assert.Equal(t, mockIPList_expected, mockIPList)

	})

	t.Run("Test with invalid JSON file", func(t *testing.T) {
		os.Setenv(constants.DBPATH, "random/random")
		defer os.Unsetenv(constants.DBPATH)
		err := LoadConfigIPConfiguration(mockConfig, mockIPList)
		assert.Error(t, err)
	})
	t.Run("Test with empty env variable PATH", func(t *testing.T) {
		os.Setenv(constants.DBPATH, "")
		defer os.Unsetenv(constants.DBPATH)
		err := LoadConfigIPConfiguration(mockConfig, mockIPList)
		assert.NoError(t, err)
	})
}

func TestSearch(t *testing.T) {
	l := make([]*ipconfig.IPConfigData, 0)
	l = append(l, &ipconfig.IPConfigData{
		Hostname:    "dummy_1",
		IPAddresses: "127.0.0.1",
		Status:      true,
	}, &ipconfig.IPConfigData{
		Hostname:    "dummy_2",
		IPAddresses: "127.0.0.2",
		Status:      false,
	},
		&ipconfig.IPConfigData{
			Hostname:    "dummy_3",
			IPAddresses: "127.0.0.3",
			Status:      false,
		},
		&ipconfig.IPConfigData{
			Hostname:    "dummy_3",
			IPAddresses: "127.0.0.4",
			Status:      true,
		},
		&ipconfig.IPConfigData{
			Hostname:    "dummy_3",
			IPAddresses: "127.0.0.5",
			Status:      true,
		})
	t.Run("Positive Test for Binary Search", func(t *testing.T) {
		index := Search(l, "127.0.0.5")
		assert.Equal(t, 4, index)
	})
	t.Run("Negative Test for Binary Search", func(t *testing.T) {
		index := Search(l, "127.0.0.9")
		assert.Equal(t, -1, index)
	})
}

func TestTTLForFileSaving(t *testing.T) {
	mockMap := ipconfig.NewMap()
	mockIPList_expected := ipconfig.NewIPConfigList()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	os.Setenv(constants.DBPATH, "/tmp/code/mock/test_data/ipconfig_test.json")
	defer os.Unsetenv(constants.DBPATH)
	// Initialize a mock IPListI
	mocklist := ipconfig.NewIPConfigList()
	l := make([]*ipconfig.IPConfigData, 0)
	l = append(l, &ipconfig.IPConfigData{
		Hostname:    "dummy_1",
		IPAddresses: "127.0.0.1",
		Status:      true,
	}, &ipconfig.IPConfigData{
		Hostname:    "dummy_2",
		IPAddresses: "127.0.0.2",
		Status:      false,
	},
		&ipconfig.IPConfigData{
			Hostname:    "dummy_3",
			IPAddresses: "127.0.0.3",
			Status:      false,
		},
		&ipconfig.IPConfigData{
			Hostname:    "dummy_3",
			IPAddresses: "127.0.0.4",
			Status:      true,
		},
		&ipconfig.IPConfigData{
			Hostname:    "dummy_3",
			IPAddresses: "127.0.0.5",
			Status:      true,
		})
	mocklist.SetIPList(l)

	Ticker = time.NewTicker(100 * time.Millisecond)
	defer Ticker.Stop()
	FLAGTOSAVE = true
	natsConn := new(mocking.MockNATSConn)
	s := "Roll Back"
	bye, _ := json.Marshal(s)
	natsConn.On("Publish", constants.CONFIGSERVICE_PUB_SUBJECT, bye).Return(nil)
	// Execute the TTLForFileSaving function in a goroutine
	go TTLForFileSaving(ctx, mocklist, natsConn)

	// Wait for a while to allow the function to execute
	time.Sleep(500 * time.Millisecond)
	LoadConfigIPConfiguration(mockMap, mockIPList_expected)
	// fmt.Println(mockIPList_expected.GetIPValues())
	// fmt.Println(mocklist.GetIPValues())
	sort.Slice(mockIPList_expected.GetIPValues(), func(i, j int) bool {
		return mockIPList_expected.GetIPValues()[i].Hostname < mockIPList_expected.GetIPValues()[j].Hostname
	})
	sort.Slice(mocklist.GetIPValues(), func(i, j int) bool {
		return mocklist.GetIPValues()[i].Hostname < mocklist.GetIPValues()[j].Hostname
	})

	assert.Equal(t, mockIPList_expected, mocklist)

	// Add assertions and validations based on your specific requirements
}
