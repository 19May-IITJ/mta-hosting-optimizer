package loader

import (
	"mta2/modules/configservice/cinternals/constants"
	"mta2/modules/configservice/cpkg/ipconfig"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfigIPConfiguration(t *testing.T) {
	// Create mock ipconfig.Configuration and ipconfig.IPList objects
	mockConfig := ipconfig.NewMap()
	mockIPList := ipconfig.NewIPConfigList()

	t.Run("Test with valid JSON file path", func(t *testing.T) {
		os.Setenv(constants.DBPATH, "/Users/b0268986/mta2/mock/test_data/ipconfig_test.json")
		defer os.Unsetenv(constants.DBPATH)
		err := LoadConfigIPConfiguration(mockConfig, mockIPList)
		assert.NoError(t, err)
		mockConfig_expected := ipconfig.NewMap()
		mockIPList_expected := ipconfig.NewIPConfigList()
		mockConfig_expected.Put("dummy_1", &ipconfig.HostData{
			HostedIP: []string{"127.0.0.1-0", "127.0.0.2-0", "127.0.0.3-1"},
			ActiveIP: 1,
		})
		mockConfig_expected.Put("dummy_2", &ipconfig.HostData{
			HostedIP: []string{"127.0.0.4-1", "127.0.0.5-1"},
			ActiveIP: 2,
		})
		mockConfig_expected.Put("dummy_3", &ipconfig.HostData{
			HostedIP: []string{"127.0.0.6-0"},
			ActiveIP: 0,
		})

		l := make([]*ipconfig.IPConfigData, 0)
		l = append(l,
			&ipconfig.IPConfigData{
				Hostname:    "dummy_1",
				IPAddresses: "127.0.0.1",
				Status:      false,
			},
			&ipconfig.IPConfigData{
				Hostname:    "dummy_1",
				IPAddresses: "127.0.0.2",
				Status:      false,
			},
			&ipconfig.IPConfigData{
				Hostname:    "dummy_1",
				IPAddresses: "127.0.0.3",
				Status:      true,
			},
			&ipconfig.IPConfigData{
				Hostname:    "dummy_2",
				IPAddresses: "127.0.0.4",
				Status:      true,
			},
			&ipconfig.IPConfigData{
				Hostname:    "dummy_2",
				IPAddresses: "127.0.0.5",
				Status:      true,
			},
			&ipconfig.IPConfigData{
				Hostname:    "dummy_3",
				IPAddresses: "127.0.0.6",
				Status:      false,
			},
		)
		mockIPList_expected.SetIPList(l)
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
